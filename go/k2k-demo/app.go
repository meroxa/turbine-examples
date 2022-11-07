package main

import (

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
	cck, err := v.Resources("cck")
	if err != nil {
		return err
	}

	rr, err := cck.Records("topic_1", []turbine.ResourceConfig{
		{"conduit", "true"},
	})

	if err != nil {
		return err
	}

	// collection doesn't work for Kafka destination
	err = cck.WriteWithConfig(rr, "topic_2", []turbine.ResourceConfig{
		{"conduit", "true"},
		{"topic", "topic_2"},
	})
	if err != nil {
		return err
	}

	return nil
}
