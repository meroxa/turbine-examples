package main

import (
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

func (a App) Run(v turbine.Turbine) error {
	source, err := v.Resources("cck")
	if err != nil {
		return err
	}

	records, err := source.Records("inbound", nil)
	if err != nil {
		return err
	}

	processed, err := v.Process(records, Format{})
	if err != nil {
		return err
	}

	dest, err := v.Resources("demopg")
	if err != nil {
		return err
	}

	err = dest.WriteWithConfig(processed, "inbound_events", turbine.ConnectionOptions{
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
