package main

import (
	"fmt"
	"log"
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

type Event struct {
	eventType                  string
	eventMultiplier            string
	players                    int
	participationPoints        int
	format                     string
	place                      uint
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
		place:                      0,
		sanctioningNumber:          "",
		matches:                    matches,
		planeswalkerPointsYearly:   0,
		planeswalkerPointsLifetime: 0,
	}
}

func _printEvent(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode || n.Type == html.TextNode {
		if n.Data != "" && n.Type == html.TextNode {
			fmt.Printf("Node value: %s\n", n.Data)
		}

		// if n.Data == "tr" {
		// 	fmt.Println("hei")
		// }

		for _, a := range n.Attr {
			fmt.Printf("Key: %s, Value: %s\n", a.Key, a.Val)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = _printEvent(links, c)
	}
	return links
}

func printEvent(eventData string) {
	doc, err := html.Parse(strings.NewReader(eventData))
	if err != nil {
		log.Fatal(err)
	}

	var parsedResults []string
	_printEvent(parsedResults, doc)
}

func parseOpponents(n *html.Node) []string {
	var opponents []string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		for _, a := range c.Attr {
			if a.Key == "class" && a.Val == "TeamOpponent" {
				opponent := c.FirstChild.Data
				opponents = append(opponents, opponent)
			}
		}
	}

	return opponents
}

func parseMatch(n *html.Node) Match {
	var place string
	var result MatchResult
	var points string
	var opponents []string

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		fmt.Printf("Data: %s, Namespace: %s, Attr: %v\n", c.Data, c.Namespace, c.Attr)

		for _, a := range c.Attr {
			if a.Key == "class" && a.Val == "MatchPlace" {
				place = c.FirstChild.Data
				fmt.Printf("Match place: %s\n", place)
			}
			if a.Key == "class" && a.Val == "MatchResult" {
				result := NewMatchResult(c.FirstChild.Data)
				fmt.Printf("Match result: %s\n", result)
			}
			if a.Key == "class" && a.Val == "MatchPoints" {
				points = c.FirstChild.Data
				fmt.Printf("Match Points: %s\n", points)
			}

			if a.Key == "class" && a.Val == "MatchOpponent" {
				opponents = parseOpponents(c)
				fmt.Printf("Match opponent: %v\n", opponents)
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

func _parseEvent(parsedEvent *Event, n *html.Node) {

	// if strings.TrimSpace(n.Data) == "MatchHistoryRow" || true {
	// 	trimmedData := strings.TrimSpace((n.Data))
	// 	if trimmedData != "" {
	// 		fmt.Printf("Node value: %s, nodetype: %v\n", trimmedData, n.Type)
	// 	}
	// }

	if n.Type == html.ElementNode && n.Data == "tr" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "MatchHistoryRow" {
				_ = parseMatch(n)
			}
		}
	}

	// if n.Type == html.ElementNode || n.Type == html.TextNode {
	// 	if n.Data != "" && n.Type == html.TextNode {
	// 		fmt.Printf("Node value: %s\n", n.Data)
	// 	}

	// 	for _, a := range n.Attr {
	// 		fmt.Printf("Key: %s, Value: %s\n", a.Key, a.Val)
	// 	}
	// }

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
}
