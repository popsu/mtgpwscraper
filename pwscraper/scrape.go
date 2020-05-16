package pwscraper

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

const (
	saveFolder           = "data"
	saveFileFormat       = "event-%s.html"
	saveFileEventHistory = "eventnames.html"
	saveJSONFilename     = "planeswalkerhistory.json"
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
	log.Printf("Saving Eventhistory to %s\n", filename)
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

func parseEventFile(filename string) *EventDetails {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	return parseEvent(text)
}

func parseHistoryFile(filename string) *EventHistory {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	text := string(content)

	return parseHistory(text)
}

func parseAllHistoryFiles(dciNumber, datafolder, saveFile string) {
	files, err := ioutil.ReadDir(datafolder)
	if err != nil {
		log.Fatal(err)
	}

	eventHistory := parseHistoryFile(path.Join(datafolder, saveFileEventHistory))
	events := make(map[string]EventDetails)
	r := regexp.MustCompile(`^event-([0-9]*)[.]html$`)

	// Parse
	for _, file := range files {
		matches := r.FindStringSubmatch((file.Name()))
		if len(matches) > 1 {
			eventid := matches[1]
			parsedEvent := parseEventFile(path.Join(datafolder, file.Name()))
			events[eventid] = *parsedEvent
		}
	}

	// Combine
	allEvents := AllEvents{DCINumber: dciNumber}

	keys := make([]string, 0, len(events))

	// Sort by numerical value instead of string value
	for k := range events {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		inum, err := strconv.Atoi(keys[i])
		if err != nil {
			// return false
			log.Fatal(err)
		}
		jnum, err := strconv.Atoi(keys[j])
		if err != nil {
			// return false
			log.Fatal(err)
		}
		return inum < jnum
	})

	for _, k := range keys {
		if _, ok := events[k]; ok {
			fullEvent := FullEvent{
				EventInfo:    eventHistory.events[k],
				EventDetails: events[k],
			}

			allEvents.AddEvent(&fullEvent)
		}
	}

	// Write to file
	log.Printf("Parsing all html into %s\n", saveFile)
	allEvents.toJson(saveFile)
}

func fetchAndSavePointHistory() string {
	pointHistory := getPointHistory(dciNumber)
	savePointHistory(pointHistory)
	return pointHistory
}

func FetcAndSaveAllEvents(pointHistory, historyFilename string) {
	eventIDs := parseEventIDs(pointHistory)
	log.Printf("Fetching %d events in parallel\n", len(eventIDs))

	var wg sync.WaitGroup
	for _, eventID := range eventIDs {
		wg.Add(1)
		go fetchAndSaveEventData(eventID, &wg)
	}
	wg.Wait()
}

func createSaveFolder() {
	err := os.Mkdir(saveFolder, 0700)
	if err != nil {
		log.Fatalf("Save folder %q already exists. Please remove it and retry. %s", saveFolder, err)
	}
}

func Execute() {
	parseFlags()
	// Scrape
	createSaveFolder()
	pointHistory := fetchAndSavePointHistory()
	historyFilename := path.Join(saveFolder, saveFileEventHistory)
	FetcAndSaveAllEvents(pointHistory, historyFilename)

	// Parse all saved html files into one JSON
	parseAllHistoryFiles(dciNumber, saveFolder, saveJSONFilename)
}
