package main

import (
	"fmt"
	"log"
	"os"

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
	source, err := v.Resources("discord")
	if err != nil {
		return err
	}

	// collection is ignored by HTTP Source Connector
	rr, err := source.Records("", turbine.ResourceConfigs{
		{"http.request.headers", authHeader()},
		{"http.request.url", guildMembersURL()}, // limit, after (highest ID)
		{"http.offset.initial", "id=0"},
		{"http.request.params", "after=${offset.id}"},
		{"http.response.record.offset.pointer", "key=/user/id, id=/user/id"},
		{"http.response.list.pointer", "/"},
	})
	if err != nil {
		return err
	}

	res := v.Process(rr, Logger{})

	// Identify a downstream data store for your data app
	// with the `Resources` function
	// Replace `destination_name` with the resource name the
	// data store was configured with on Meroxa.

	dest, err := v.Resources("s3")
	if err != nil {
		return err
	}

	// Specify where to write records downstream
	// using the `Write` function
	// Replace `collection_archive` with a table, collection,
	// or bucket name in your data store.
	// If a configuration is needed, you can also use i.e.
	//
	// dest.WriteWithConfig(
	//  res,
	//  "my-archive",
	//  turbine.ResourceConfigs{turbine.ResourceConfig{Field: "buffer.flush.time", Value: "10"}}
	// )

	err = dest.Write(res, "discord_users")
	if err != nil {
		return err
	}

	return nil
}

type Logger struct{}

func (f Logger) Process(stream []turbine.Record) []turbine.Record {
	for _, record := range stream {
		log.Printf("Record: %v", record)
	}
	return stream
}

func authHeader() string {
	authToken := os.Getenv("DISCORD_BOT_TOKEN")
	return fmt.Sprintf("Authorization: Bot %s", authToken)
}

func guildMembersURL() string {
	guildID := os.Getenv("GUILD_ID")
	return fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members", guildID)
}
