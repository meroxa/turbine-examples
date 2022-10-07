package main

import (
	// Dependencies of the example data app
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	// Dependencies of Turbine
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	// To configure your data stores as resources on the Meroxa Platform
	// use the Meroxa Dashboard, CLI, or Meroxa Terraform Provider
	// For more details refer to: https://docs.meroxa.com/

	// Identify an upstream data store for your data app
	// with the `Resources` function
	// Replace `source_name` with the resource name the
	// data store was configured with on Meroxa
	source, err := v.Resources("mdb")
	if err != nil {
		return err
	}

	// Specify which upstream records to pull
	// with the `Records` function
	// Replace `collection_name` with a table, collection,
	// or bucket name in your data store
	rr, err := source.Records("events", nil)
	if err != nil {
		return err
	}

	// Specify what code to execute against upstream records
	// with the `Process` function
	// Replace `Anonymize` with the name of your function code
	res, _ := v.Process(rr, Anonymize{})

	// Identify a downstream data store for your data app
	// with the `Resources` function
	// Replace `destination_name` with the resource name the
	// data store was configured with on Meroxa
	dest, err := v.Resources("demopg")
	if err != nil {
		return err
	}

	// Specify where to write records downstream
	// using the `Write` function
	// Replace `collection_archive` with a table, collection,
	// or bucket name in your data store
	err = dest.Write(res, "mdb_events", nil)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(stream []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	for i, r := range stream {
		rec, err := parseCDCRecord(r.Payload)
		if err != nil {
			log.Printf("error: %s", err.Error())
			break
		}

		var after map[string]interface{}
		err = json.Unmarshal([]byte(rec.Payload.After), &after)
		if err != nil {
			log.Printf("error: %s", err.Error())
			break
		}

		after["email"] = consistentHash(after["email"].(string))

		afterJSON, err := json.Marshal(after)
		if err != nil {
			log.Printf("error: %s", err.Error())
			break
		}

		rec.Payload.After = string(afterJSON)
		b, err := rec.Bytes()
		if err != nil {
			log.Printf("error: %s", err.Error())
			break
		}
		r.Payload = b
		stream[i] = r
	}

	return stream, nil
}

func parseCDCRecord(raw []byte) (CDCRecord, error) {
	var rec CDCRecord
	err := json.Unmarshal(raw, &rec)
	if err != nil {
		return CDCRecord{}, err
	}
	return rec, nil
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
