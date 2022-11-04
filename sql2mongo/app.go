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

	source, err := v.Resources("sqlpg")
	if err != nil {
		return err
	}

	query := `SELECT * FROM 
             (SELECT purchases.id, purchases.created_at, users.country, users.city
              FROM purchases
              	INNER JOIN users
              	ON user_id = users.id) as a
`

	rr, err := source.Records("", turbine.ConnectionOptions{
		turbine.ConnectionOption{
			Field: "mode",
			Value: "timestamp+incrementing",
		},
		turbine.ConnectionOption{
			Field: "incrementing.column.name",
			Value: "id",
		},
		turbine.ConnectionOption{
			Field: "timestamp.column.name",
			Value: "created_at",
		},
		turbine.ConnectionOption{
			Field: "query",
			Value: query,
		},
	})
	if err != nil {
		return err
	}

	res := v.Process(rr, Transform{})

	dest, err := v.Resources("mdb")
	if err != nil {
		return err
	}

	err = dest.Write(res, "collection_archive")
	if err != nil {
		return err
	}

	return nil
}

type Transform struct{}

func (f Transform) Process(stream []turbine.Record) []turbine.Record {
	for i, record := range stream {
		log.Printf("record: %+v", record)
		record.Payload.Set("migration", "true")
		stream[i] = record
	}
	return stream
}
