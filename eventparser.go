package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type MatchResult int

const (
	Bye MatchResult = iota
	Win
	Loss
	Draw
)

func (m MatchResult) String() string {
	return [...]string{"Bye", "Win", "Loss", "Draw"}[m]
}

func (m MatchResult) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(m.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func NewMatchResult(result string) MatchResult {
	switch result {
	case "Bye":
		return Bye
	case "Win":
		return Win
	case "Loss":
		return Loss
	case "Draw":
		return Draw
	}
	panic("Invalid matchresult")
}

type Match struct {
	Place     string      `json:",omitempty"`
	Result    MatchResult `json:",omitempty"`
	Points    string      `json:",omitempty"`
	Opponents []string    `json:",omitempty"`
}

func (m Match) String() string {
	return fmt.Sprintf("Place: %s, Result %s, Points: %s, Opponents: %s",
		m.Place, m.Result, m.Points, m.Opponents)
}

type Event struct {
	EventType                      string  `json:",omitempty"`
	EventMultiplier                string  `json:",omitempty"`
	Players                        int     `json:",omitempty"`
	ParticipationPoints            int     `json:",omitempty"`
	Format                         string  `json:",omitempty"`
	Location                       string  `json:",omitempty"`
	Place                          int     `json:",omitempty"`
	SanctioningNumber              string  `json:",omitempty"`
	Matches                        []Match `json:",omitempty"`
	PlaneswalkerPointsYearly       int     `json:",omitempty"`
	PlaneswalkerPointsProfessional int     `json:",omitempty"`
	PlaneswalkerPointsLifetime     int     `json:",omitempty"`
}

func (e Event) toJson() {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("testdata.json", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// addMatch adds a Match to an Event
func (e *Event) addMatch(m Match) {
	e.Matches = append(e.Matches, m)
}

func (e *Event) addPlaneswalkerPoints(n *html.Node) {
	switch n.FirstChild.NextSibling.FirstChild.Data {
	case "Yearly:":
		e.PlaneswalkerPointsYearly = _parseEventInt(n)
	case "Professional:":
		e.PlaneswalkerPointsProfessional = _parseEventInt(n)
	case "Lifetime:":
		e.PlaneswalkerPointsLifetime = _parseEventInt(n)
	}
}

func (e Event) String() string {
	var s string

	s = s + fmt.Sprintf("EventType           : %s\n", e.EventType)
	s = s + fmt.Sprintf("EventMultiplier     : %s\n", e.EventMultiplier)
	s = s + fmt.Sprintf("EventPlayers        : %d\n", e.Players)
	s = s + fmt.Sprintf("Participation points: %d\n", e.ParticipationPoints)
	s = s + fmt.Sprintf("Format              : %s\n", e.Format)
	s = s + fmt.Sprintf("Location:           : %s\n", e.Location)
	s = s + fmt.Sprintf("Place:              : %d\n", e.Place)
	s = s + fmt.Sprintf("Sanctioning number  : %s\n", e.SanctioningNumber)
	s = s + fmt.Sprintf("Planeswalker Points\n")
	s = s + fmt.Sprintf("    Yearly          : %d\n", e.PlaneswalkerPointsYearly)
	s = s + fmt.Sprintf("    Professional    : %d\n", e.PlaneswalkerPointsProfessional)
	s = s + fmt.Sprintf("    Lifetime        : %d\n", e.PlaneswalkerPointsLifetime)

	for _, m := range e.Matches {
		s = s + fmt.Sprintf("Match: %s\n", m)
	}

	return s
}

func _parseOpponents(opponents *[]string, n *html.Node) {
	for _, a := range n.Attr {
		if a.Key == "class" && a.Val == "TeamOpponent" {
			opponent := n.FirstChild.Data
			*opponents = append(*opponents, opponent)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		_parseOpponents(opponents, c)
	}
}

func parseOpponents(n *html.Node) []string {
	var opponents []string
	_parseOpponents(&opponents, n)

	return opponents
}

func parseMatch(n *html.Node) Match {
	var place string
	var result MatchResult
	var points string
	var opponents []string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Key == "class" && a.Val == "MatchPlace" {
				place = c.FirstChild.Data
			}
			if a.Key == "class" && a.Val == "MatchResult" {
				result = NewMatchResult(c.FirstChild.Data)
			}
			if a.Key == "class" && a.Val == "MatchPoints" {
				points = c.FirstChild.Data
			}
			if a.Key == "class" && a.Val == "MatchOpponent" {
				opponents = parseOpponents(c)
			}
		}
	}
	return Match{
		Place:     place,
		Result:    result,
		Points:    points,
		Opponents: opponents,
	}
}

func _parseEventString(n *html.Node) string {
	return strings.TrimSpace(n.LastChild.Data)
}

func _parseEventInt(n *html.Node) int {
	val, err := strconv.Atoi(strings.TrimSpace(n.LastChild.Data))
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func _parseEvent(parsedEvent *Event, n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" {
				if a.Val == "EventType" {
					parsedEvent.EventType = _parseEventString(n)
				} else if a.Val == "EventMultiplier" {
					parsedEvent.EventMultiplier = _parseEventString(n)
				} else if a.Val == "EventPlayers" {
					parsedEvent.Players = _parseEventInt(n)
				} else if a.Val == "EventParticipationPoints" {
					parsedEvent.ParticipationPoints = _parseEventInt(n)
				} else if a.Val == "EventFormat" {
					parsedEvent.Format = _parseEventString(n)
				} else if a.Val == "EventLocation" {
					parsedEvent.Location = _parseEventString(n)
				} else if a.Val == "EventPlace" {
					parsedEvent.Place = _parseEventInt(n)
				} else if a.Val == "EventSanctionNumber" {
					parsedEvent.SanctioningNumber = _parseEventString(n)
				} else if a.Val == "MatchTotal" {
					parsedEvent.addPlaneswalkerPoints(n)
				}
			}
		}
	}

	if n.Type == html.ElementNode && n.Data == "tr" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "MatchHistoryRow" {
				parsedEvent.addMatch(parseMatch(n))
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		_parseEvent(parsedEvent, c)
	}

}

func parseEvent(eventData string) *Event {
	doc, err := html.Parse(strings.NewReader(eventData))
	if err != nil {
		log.Fatal(err)
	}

	var parsedEvent = &Event{}
	_parseEvent(parsedEvent, doc)

	return parsedEvent
}
