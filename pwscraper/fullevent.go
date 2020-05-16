package pwscraper

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type FullEvent struct {
	EventInfo    EventInfo
	EventDetails EventDetails
}

func (fe FullEvent) toJson(filename string) {
	bytes, err := json.Marshal(fe)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

type AllEvents struct {
	Events    []FullEvent
	DCINumber string
}

func (ae *AllEvents) AddEvent(fe *FullEvent) {
	ae.Events = append(ae.Events, *fe)
}

func (ae AllEvents) toJson(filename string) {
	bytes, err := json.MarshalIndent(ae, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
