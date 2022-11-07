package main

import (
	"encoding/json"
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
	"log"
	"strconv"
	"time"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {

	// reference the MongoDB resource that was created on the platform. In this case I created "mdb".
	source, err := v.Resources("mdb")
	if err != nil {
		return err
	}

	// pull records from the "events" collection.
	rr, err := source.Records("events", nil)
	if err != nil {
		return err
	}

	// apply the "FilterInteresting" processor to those records.
	res := v.Process(rr, FilterInteresting{})

	// reference the Kafka resource that was created on the platform. In this case I created "cck".
	dest, err := v.Resources("cck")
	if err != nil {
		return err
	}

	// write out the resulting records into the collection (or __Topic__ in the case of Kafka). In this case I'm writing
	// out to the Topic "interesting_events".
	err = dest.WriteWithConfig(res, "interesting_events", nil)
	if err != nil {
		return err
	}

	return nil
}

// FilterInteresting looks for "interesting" events and filters out everything else. For this example, __interesting__
// events are any events where an event is associated with a VIP user.
type FilterInteresting struct{}

func (f FilterInteresting) Process(stream []turbine.Record) []turbine.Record {
	var interestingEvents []Event
	for _, r := range stream {
		ev, err := parseEventRecord(r)
		if err != nil {
			log.Printf("error: %s", err.Error())
			break
		}

		if isInteresting(ev) {
			interestingEvents = append(interestingEvents, ev)
		}
	}

	if len(interestingEvents) > 0 {
		recs, err := encodeEvents(interestingEvents)
		if err != nil {
			log.Printf("error: %s", err.Error())
		}
		return recs
	}
	return []turbine.Record{}
}

// Event represents the Event document stored in MongoDB.
type Event struct {
	UserID    string    `json:"user_id"`
	Activity  string    `json:"activity"`
	VIP       bool      `json:"vip"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

// UnmarshalJSON is a custom unmarshaler for Event, that handles the conversion of the VIP field from string to bool
func (ev *Event) UnmarshalJSON(data []byte) error {
	type EvAlias Event
	aux := &struct {
		*EvAlias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		DeletedAt string `json:"deleted_at"`
	}{
		EvAlias: (*EvAlias)(ev),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	iCreatedAt, err := strconv.Atoi(aux.CreatedAt)
	if err == nil {
		ev.CreatedAt = time.Unix(int64(iCreatedAt), 0)
	}
	iUpdatedAt, err := strconv.Atoi(aux.UpdatedAt)
	if err == nil {
		ev.UpdatedAt = time.Unix(int64(iUpdatedAt), 0)
	}
	iDeletedAt, err := strconv.Atoi(aux.DeletedAt)
	if err == nil {
		ev.DeletedAt = time.Unix(int64(iDeletedAt), 0)
	}
	return nil
}

// parses a turbine.Record and returns an Event struct
func parseEventRecord(r turbine.Record) (Event, error) {
	// todo: this doesn't work. We need to pull "after" out of the payload
	var cdcRec CDCRecord
	err := json.Unmarshal(r.Payload, &cdcRec)
	if err != nil {
		return Event{}, err
	}
	log.Printf("After: %s", cdcRec.Payload.After)
	var ev Event
	err = json.Unmarshal([]byte(cdcRec.Payload.After), &ev)
	if err != nil {
		return Event{}, err
	}
	log.Printf("Event Record: %+v", ev)
	return ev, nil
}

// encode a slice of Events into a slice of turbine.Record
func encodeEvents(ee []Event) ([]turbine.Record, error) {
	var rr []turbine.Record
	for _, ev := range ee {
		b, err := json.Marshal(ev)
		if err != nil {
			return nil, err
		}
		rr = append(rr, turbine.Record{
			Key:       ev.UserID,
			Payload:   b,
			Timestamp: ev.UpdatedAt,
		})
	}
	log.Printf("Emitted record: %+v", rr)
	return rr, nil
}

// in this case it's as simple as returning the value ev.VIP since it defaults to false, but
// we could do anything here. E.g. apply some function on various fields or check with an
// external API/service.
func isInteresting(ev Event) bool {
	return ev.VIP
}
