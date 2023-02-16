package main

import (
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
	for _, record := range stream {
		payload, err := record.Payload.Map()
		if err != nil {
			log.Print("unexpected payload", err)
			break
		}
		log.Printf("record payload: %+v", payload)
	}
	return stream
}
