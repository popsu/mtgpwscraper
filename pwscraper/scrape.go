package pwscraper

import (
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

type Config struct {
	DciNumber            string
	PwpCookieValue       string
	SaveFolder           string
	SaveFileFormat       string
	SaveFileEventHistory string
	SaveJSONFilename     string
}

func NewDefaultConfig() *Config {
	return &Config{
		DciNumber:            "",
		PwpCookieValue:       "",
		SaveFolder:           "data",
		SaveFileFormat:       "event-%s.html",
		SaveFileEventHistory: "eventnames.html",
		SaveJSONFilename:     "planeswalkerhistory.json",
	}
}

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

func fetchAndSaveEventData(eventID, pwpCookieValue, saveFolder, saveFileFormat string, wg *sync.WaitGroup) {
	defer wg.Done()

	data := []byte(fetchEventData(eventID, pwpCookieValue))

	filename := path.Join(saveFolder, fmt.Sprintf(saveFileFormat, eventID))
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func savePointHistory(pointHistory, saveFolder, saveFileEventHistory string) {
	filename := path.Join(saveFolder, saveFileEventHistory)
	log.Printf("Saving Eventhistory to %s\n", filename)
	err := ioutil.WriteFile(filename, []byte(pointHistory), 0644)
	if err != nil {
		log.Fatal(err)
	}
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

func parseAllHistoryFiles(dciNumber, datafolder, saveFileJSON, saveFileEventHistory string) {
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
	log.Printf("Parsing all html into %s\n", saveFileJSON)
	allEvents.toJson(saveFileJSON)
}

func fetchAndSavePointHistory(dciNumber, saveFolder, saveFileEventHistory string) string {
	pointHistory := getPointHistory(dciNumber)
	savePointHistory(pointHistory, saveFolder, saveFileEventHistory)
	return pointHistory
}

func FetcAndSaveAllEvents(pointHistory, historyFilename, pwpCookieValue, saveFolder, saveFileFormat string) {
	eventIDs := parseEventIDs(pointHistory)
	log.Printf("Fetching %d events to %s/\n", len(eventIDs), saveFolder)
	log.Printf("If the EventDetails in the JSON file is empty, the cookie was incorrect")

	var wg sync.WaitGroup
	for _, eventID := range eventIDs {
		wg.Add(1)
		go fetchAndSaveEventData(eventID, pwpCookieValue, saveFolder, saveFileFormat, &wg)
	}
	wg.Wait()
}

func createSaveFolder(saveFolder string) {
	err := os.Mkdir(saveFolder, 0700)
	if err != nil {
		log.Fatalf("Savefolder %q already exists. Please remove it and retry. %s", saveFolder, err)
	}
}

func Execute(conf *Config) {
	// Scrape
	createSaveFolder(conf.SaveFolder)
	pointHistory := fetchAndSavePointHistory(conf.DciNumber, conf.SaveFolder, conf.SaveFileEventHistory)
	historyFilename := path.Join(conf.SaveFolder, conf.SaveFileEventHistory)
	FetcAndSaveAllEvents(pointHistory, historyFilename, conf.PwpCookieValue, conf.SaveFolder, conf.SaveFileFormat)

	// Parse all saved html files into one JSON
	parseAllHistoryFiles(conf.DciNumber, conf.SaveFolder, conf.SaveJSONFilename, conf.SaveFileEventHistory)
}
