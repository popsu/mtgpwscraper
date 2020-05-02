package main

import (
	"fmt"
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
	place    string
	result   MatchResult
	points   string
	opponent []string
}

func (m Match) String() string {
	return fmt.Sprintf("Place: %s, Result %s, Points: %s, Opponents: %s",
		m.place, m.result, m.points, m.opponent)
}

type Event struct {
	eventType                  string
	eventMultiplier            string
	players                    int
	participationPoints        int
	format                     string
	location                   string
	place                      int
	sanctioningNumber          string
	matches                    []Match
	planeswalkerPointsYearly   int
	planeswalkerPointsLifetime int
}

func NewEvent() *Event {
	var matches []Match

	return &Event{
		eventType:                  "",
		eventMultiplier:            "",
		players:                    0,
		participationPoints:        0,
		format:                     "",
		location:                   "",
		place:                      0,
		sanctioningNumber:          "",
		matches:                    matches,
		planeswalkerPointsYearly:   0,
		planeswalkerPointsLifetime: 0,
	}
}

// addMatch adds a Match to an Event
func (e *Event) addMatch(m Match) {
	e.matches = append(e.matches, m)
}

func (e *Event) addPlaneswalkerPoints(n *html.Node) {
	switch n.FirstChild.NextSibling.FirstChild.Data {
	case "Yearly:":
		e.planeswalkerPointsYearly = _parseEventInt(n)
	case "Lifetime:":
		e.planeswalkerPointsLifetime = _parseEventInt(n)
	}
}

func (e Event) String() string {
	var s string

	s = s + fmt.Sprintf("EventType           : %s\n", e.eventType)
	s = s + fmt.Sprintf("EventMultiplier     : %s\n", e.eventMultiplier)
	s = s + fmt.Sprintf("EventPlayers        : %d\n", e.players)
	s = s + fmt.Sprintf("Participation points: %d\n", e.participationPoints)
	s = s + fmt.Sprintf("Format              : %s\n", e.format)
	s = s + fmt.Sprintf("Location:           : %s\n", e.location)
	s = s + fmt.Sprintf("Place:              : %d\n", e.place)
	s = s + fmt.Sprintf("Sanctioning number  : %s\n", e.sanctioningNumber)
	s = s + fmt.Sprintf("Planeswalke Points\n")
	s = s + fmt.Sprintf("    Yearly          : %d\n", e.planeswalkerPointsYearly)
	s = s + fmt.Sprintf("    Lifetime        : %d\n", e.planeswalkerPointsLifetime)

	for _, m := range e.matches {
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
		place:    place,
		result:   result,
		points:   points,
		opponent: opponents,
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
					parsedEvent.eventType = _parseEventString(n)
				} else if a.Val == "EventMultiplier" {
					parsedEvent.eventMultiplier = _parseEventString(n)
				} else if a.Val == "EventPlayers" {
					parsedEvent.players = _parseEventInt(n)
				} else if a.Val == "EventParticipationPoints" {
					parsedEvent.participationPoints = _parseEventInt(n)
				} else if a.Val == "EventFormat" {
					parsedEvent.format = _parseEventString(n)
				} else if a.Val == "EventLocation" {
					parsedEvent.location = _parseEventString(n)
				} else if a.Val == "EventPlace" {
					parsedEvent.place = _parseEventInt(n)
				} else if a.Val == "EventSanctionNumber" {
					parsedEvent.sanctioningNumber = _parseEventString(n)
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

func parseEvent(eventData string) {
	doc, err := html.Parse(strings.NewReader(eventData))
	if err != nil {
		log.Fatal(err)
	}

	var parsedEvent = NewEvent()
	_parseEvent(parsedEvent, doc)
	fmt.Println(parsedEvent)
}
