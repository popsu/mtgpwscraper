package main

import (
	"fmt"
	"log"
	"strings"

	"cloud.google.com/go/civil"
	"golang.org/x/net/html"
)

type EventHistory struct {
	events map[string]EventInfo
}

func NewEventHistory() *EventHistory {
	events := make(map[string]EventInfo)
	return &EventHistory{
		events: events,
	}
}

func (e EventHistory) String() string {
	var s string

	for _, v := range e.events {
		s = s + fmt.Sprintf("%s\n", v)
	}
	return s
}

type EventInfo struct {
	date           civil.Date
	description    string
	location       string
	lifeTimePoints string
	proPoints      string
	ID             string
}

func (e EventInfo) String() string {
	var s string

	s = s + fmt.Sprintf("Date           : %s\n", e.date)
	s = s + fmt.Sprintf("Description    : %s\n", e.description)
	s = s + fmt.Sprintf("Location       : %s\n", e.location)
	s = s + fmt.Sprintf("LifetimePoints : %s\n", e.lifeTimePoints)
	s = s + fmt.Sprintf("ProPoints      : %s\n", e.proPoints)
	s = s + fmt.Sprintf("ID             : %s\n", e.ID)

	return s
}

func _parseEventID(n *html.Node) string {
	for _, a := range n.LastChild.PrevSibling.Attr {
		if a.Key == "data-summarykey" {
			return strings.TrimSpace(a.Val)
		}
	}

	log.Fatal("Missing eventID")
	return ""
}

func _parseEventDate(n *html.Node) civil.Date {
	dateRFC3339 := n.FirstChild.Data

	d, err := civil.ParseDate(dateRFC3339)
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func _parseEventLocation(n *html.Node) string {
	return n.FirstChild.Data
}

func _parseEventDescription(n *html.Node) string {
	return n.FirstChild.Data
}

func _parseEventLifetimePoints(n *html.Node) string {
	return n.FirstChild.Data
}

func _parseEventProPoints(n *html.Node) string {
	return n.FirstChild.Data
}

func parseHistoryEvent(eventInfo *EventInfo, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "UnLocked" {
				eventInfo.ID = _parseEventID(n)
			} else if a.Key == "class" && a.Val == "HistoryPanelHeaderLabel Date" {
				eventInfo.date = _parseEventDate(n)
			} else if a.Key == "class" && a.Val == "HistoryPanelHeaderLabel Description" {
				eventInfo.description = _parseEventDescription(n)
			} else if a.Key == "class" && a.Val == "HistoryPanelHeaderLabel Location" {
				eventInfo.location = _parseEventLocation(n)
			} else if a.Key == "class" && a.Val == "HistoryPanelHeaderLabel LifetimePoints" {
				eventInfo.lifeTimePoints = _parseEventLifetimePoints(n)
			} else if a.Key == "class" && a.Val == "HistoryPanelHeaderLabel ProPoints" {
				eventInfo.proPoints = _parseEventProPoints(n)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseHistoryEvent(eventInfo, c)
	}

}

func _parseHistory(parsedHistory *EventHistory, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "HistoryPanelRow" {
				eventInfo := &EventInfo{}
				parseHistoryEvent(eventInfo, n)
				parsedHistory.events[eventInfo.ID] = *eventInfo
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		_parseHistory(parsedHistory, c)
	}

}

func parseHistory(eventData string) {
	doc, err := html.Parse(strings.NewReader(eventData))
	if err != nil {
		log.Fatal(err)
	}

	parsedHistory := NewEventHistory()
	_parseHistory(parsedHistory, doc)
	fmt.Println(parsedHistory.events["896252"])
	fmt.Println(parsedHistory)
}
