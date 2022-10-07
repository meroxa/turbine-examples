package main

import (

	// Dependencies of Turbine
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
	"log"

	"github.com/ahamidi/kcschema"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	source, err := v.Resources("cck")
	if err != nil {
		return err
	}

	rr, err := source.Records("inbound", nil)
	if err != nil {
		return err
	}

	res := v.Process(rr, Format{})

	dest, err := v.Resources("demopg")
	if err != nil {
		return err
	}

	err = dest.WriteWithConfig(res, "inbound_events", turbine.ResourceConfigs{
		{Field: "key.converter", Value: "org.apache.kafka.connect.storage.StringConverter"},
		{Field: "key.converter.schemas.enable", Value: "true"},
	})
	if err != nil {
		return err
	}

	return nil
}

type Format struct{}

func (f Format) Process(stream []turbine.Record) []turbine.Record {
	for i, record := range stream {
		log.Printf("Original turbine Record: %+v", record)
		sp, err := kcschema.Parse(kcschema.Payload(record.Payload))
		if err != nil {
			log.Printf("error casting Payload to Map: %s", err.Error())
			break
		}
		log.Printf("Parsed payload: %+v", sp)

		j, err := sp.AsKCSchemaJSON("inbound")
		if err != nil {
			log.Printf("error casting Payload to Map: %s", err.Error())
			break
		}
		log.Printf("converted record with schema: %+v", string(j))
		stream[i].Payload = j
	}
	return stream
}
