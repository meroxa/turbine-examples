package main

import (
	"log"

	// Dependencies of Turbine
	"github.com/meroxa/turbine-go/pkg/turbine"
	"github.com/meroxa/turbine-go/pkg/turbine/cmd"
	"github.com/meroxa/turbine-go/pkg/turbine/transforms"
)

func main() {
	cmd.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	source, err := v.Resources("mongo")
	if err != nil {
		return err
	}

	records, err := source.Records("events", nil)
	if err != nil {
		return err
	}

	processed, err := v.Process(records, Flatten{})
	if err != nil {
		return err
	}

	dest, err := v.Resources("destination_name")
	if err != nil {
		return err
	}

	err = dest.Write(processed, "collection_archive")
	if err != nil {
		return err
	}

	return nil
}

type Flatten struct{}

func (f Flatten) Process(stream []turbine.Record) []turbine.Record {
	for i, r := range stream {
		err := transforms.Flatten(&r.Payload)
		if err != nil {
			log.Printf("error: %s", err.Error())
		}

		stream[i] = r
	}
	return stream
}
