package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
	"golang.org/x/net/html"
)

const (
	eventSummaryURL = "https://www.wizards.com/Magic/PlaneswalkerPoints/JavaScript/GetEventSummary/"
	pointHistoryURL = "https://www.wizards.com/Magic/PlaneswalkerPoints/JavaScript/GetPointsHistory/%s"

	pwpCookieNameFmt = "PWP.ASPXAUTH=%s"

	pointHistoryJSONPath = "Data.1.Value"
	eventDataJSONPath    = "Data.Value"

	saveFolder     = "data"
	saveFileFormat = "event-%s.html"
)

var (
	pwpCookieValue = ""
	dciNumber      = ""
)

func getPointHistory(dcinumber string) string {
	URL := fmt.Sprintf(pointHistoryURL, dcinumber)

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	request, err := http.NewRequest("POST", URL, strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	res := gjson.Get(string(body), pointHistoryJSONPath).String()

	return res
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

func fetchEventData(eventID string) string {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	postPayload := fmt.Sprintf("ID=%s", eventID)
	request, err := http.NewRequest("POST", eventSummaryURL, strings.NewReader(postPayload))
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("cookie", fmt.Sprintf(pwpCookieNameFmt, pwpCookieValue))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatal(body)
	}

	return gjson.Get(string(body), eventDataJSONPath).String()
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

func parseFlags() {
	flag.StringVar(&dciNumber, "dcinumber", "", "DCI Number")
	flag.StringVar(&pwpCookieValue, "cookie", "", "Cookie named PWP.ASPXAUTH in wizards.com site")

	flag.Parse()
}

func main() {
	parseFlags()

	pointHistory := getPointHistory(dciNumber)
	eventIDs := parseEventIDs(pointHistory)

	err := os.Mkdir(saveFolder, 0700)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, eventID := range eventIDs {
		wg.Add(1)
		go fetchAndSaveEventData(eventID, &wg)
	}

	wg.Wait()
}
