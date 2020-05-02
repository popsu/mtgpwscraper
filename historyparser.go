package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type EventHistory struct {
	events map[string]EventInfo
}

func NewEventHistory() *EventHistory {
	var events map[string]EventInfo
	return &EventHistory{
		events: events,
	}
}

type EventInfo struct {
	date           time.Time
	description    string
	location       string
	lifeTimePoints int
	proPoints      int
	ID             string
}

func NewEventInfo() *EventInfo {
	return &EventInfo{
		date:           time.Time{},
		description:    "",
		location:       "",
		lifeTimePoints: 0,
		proPoints:      0,
		ID:             "",
	}
}

func (e EventInfo) String() string {
	return e.ID
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

func parseHistoryEvent(eventInfo *EventInfo, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "UnLocked" {
				eventInfo.ID = _parseEventID(n)
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
				eventInfo := NewEventInfo()
				parseHistoryEvent(eventInfo, n)
				fmt.Println(eventInfo)
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
	fmt.Println(parsedHistory)
}
