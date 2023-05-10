package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	// Dependencies of Turbine
	"github.com/meroxa/turbine-go/pkg/turbine"
	"github.com/meroxa/turbine-go/pkg/turbine/cmd"
)

func main() {
	cmd.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}
type Anonymize struct{}

func (a App) Run(v turbine.Turbine) error {
	db, err := v.Resources("demopg")
	if err != nil {
		return err
	}

	records, err := db.Records("user_activity", nil) // rr is a collection of records, can't be inspected directly
	if err != nil {
		return err
	}

	processed, err := v.Process(records, Anonymize{})
	if err != nil {
		return err
	}

	s3, err := v.Resources("s3")
	if err != nil {
		return err
	}
	err = s3.Write(processed, "data-app-archive")
	if err != nil {
		return err
	}

	return nil
}

func (f Anonymize) Process(rr []turbine.Record) []turbine.Record {
	for i, r := range rr {
		hashedEmail := consistentHash(r.Payload.Get("email").(string))
		err := r.Payload.Set("email", hashedEmail)
		if err != nil {
			log.Println("error setting value: ", err)
			break
		}
		rr[i] = r
	}
	return rr
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
