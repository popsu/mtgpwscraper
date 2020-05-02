package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

const (
	saveFolder           = "data"
	saveFileFormat       = "event-%s.html"
	saveFileEventHistory = "eventnames.html"
)

var (
	pwpCookieValue = ""
	dciNumber      = ""
)

func parseEventIDs(pointHistory string) []string {
	var eventIDs []string

	doc, err := html.Parse(strings.NewReader(pointHistory))
	if err != nil {
		log.Fatal(err)
	}
	eventIDs = _parseEventIDs(eventIDs, doc)

	return eventIDs
}

func _parseEventIDs(links []string, n *html.Node) []string {
	isEvent := false
	var dataSummarykey string

	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "data-summarykey" {
				dataSummarykey = a.Val
			}
			if a.Key == "data-type" && a.Val == "Event" {
				isEvent = true
			}
		}
	}

	if isEvent {
		links = append(links, dataSummarykey)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = _parseEventIDs(links, c)
	}
	return links
}

func fetchAndSaveEventData(eventID string, wg *sync.WaitGroup) {
	defer wg.Done()

	data := []byte(fetchEventData(eventID))

	filename := path.Join(saveFolder, fmt.Sprintf(saveFileFormat, eventID))
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func savePointHistory(pointHistory string) {
	filename := path.Join(saveFolder, saveFileEventHistory)
	err := ioutil.WriteFile(filename, []byte(pointHistory), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFlags() {
	flag.StringVar(&dciNumber, "dcinumber", "", "DCI Number")
	flag.StringVar(&pwpCookieValue, "cookie", "", "Cookie named PWP.ASPXAUTH in wizards.com site")

	flag.Parse()
}

func parseEventFile(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	parseEvent(text)
}

func parseHistoryFile(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	parseHistory(text)
}

func main() {
	parseFlags()

	filename := "data/event-3596080.html"
	historyFilename := "data/eventnames.html"

	parseHistoryFile(historyFilename)

	os.Exit(0)

	parseEventFile(filename)

	pointHistory := getPointHistory(dciNumber)
	eventIDs := parseEventIDs(pointHistory)

	err := os.Mkdir(saveFolder, 0700)
	if err != nil {
		log.Fatal(err)
	}
	savePointHistory(pointHistory)

	var wg sync.WaitGroup

	for _, eventID := range eventIDs {
		wg.Add(1)
		go fetchAndSaveEventData(eventID, &wg)
	}

	wg.Wait()
}
