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
	db, err := v.Resources("demopg")
	if err != nil {
		return err
	}

	records, err := db.Records("user_activity", nil) // rr is a collection of records, can't be inspected directly
	if err != nil {
		return err
	}

	err = v.RegisterSecret("CLEARBIT_API_KEY") // makes env var available to data app
	if err != nil {
		return err
	}

	res, err := v.Process(records, EnrichUserData{})
	if err != nil {
		return err
	}

	err = db.Write(res, "user_activity_enriched")
	if err != nil {
		return err
	}

	return nil
}

type EnrichUserData struct{}

func (f EnrichUserData) Process(rr []turbine.Record) []turbine.Record {
	for i, r := range rr {
		log.Printf("Got email: %s", r.Payload.Get("email"))
		UserDetails, err := EnrichUserEmail(r.Payload.Get("email").(string))
		if err != nil {
			log.Println("error enriching user data: ", err)
			break
		}
		log.Printf("Got UserDetails: %+v", UserDetails)
		err = r.Payload.Set("full_name", UserDetails.FullName)
		if err != nil {
			log.Println("error setting full_name value: ", err)
			break
		}
		err = r.Payload.Set("company", UserDetails.Company)
		if err != nil {
			log.Println("error setting company value: ", err)
			break
		}
		err = r.Payload.Set("location", UserDetails.Location)
		if err != nil {
			log.Println("error setting location value: ", err)
			break
		}
		err = r.Payload.Set("role", UserDetails.Role)
		if err != nil {
			log.Println("error setting role value: ", err)
			break
		}
		err = r.Payload.Set("seniority", UserDetails.Seniority)
		if err != nil {
			log.Println("error setting seniority value: ", err)
			break
		}
		rr[i] = r
	}

	return rr
}
