package main

import (
	// Dependencies of the example data app
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

	source, err := v.Resources("demopg")
	if err != nil {
		return err
	}

	raw, err := source.Records("emails", nil)
	if err != nil {
		return err
	}

	res := v.Process(raw, Anonymize{})

	err = source.Write(res, "emails_processed")
	if err != nil {
		return err
	}

	dest, err := v.Resources("s3")
	if err != nil {
		return err
	}

	err = dest.Write(res, "dl-private")
	if err != nil {
		return err
	}

	err = dest.Write(raw, "dl-raw")
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(stream []turbine.Record) []turbine.Record {
	for i, record := range stream {
		email := fmt.Sprintf("%s", record.Payload.Get("email"))
		if email == "" {
			log.Printf("unable to find email value in record %d\n", i)
			break
		}
		hashedEmail := consistentHash(email)
		err := record.Payload.Set("email", hashedEmail)
		if err != nil {
			log.Println("error setting value: ", err)
			continue
		}
		stream[i] = record
	}
	return stream
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
