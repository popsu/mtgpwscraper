package pwscraper

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

const (
	eventSummaryURL = "https://www.wizards.com/Magic/PlaneswalkerPoints/JavaScript/GetEventSummary/"
	pointHistoryURL = "https://www.wizards.com/Magic/PlaneswalkerPoints/JavaScript/GetPointsHistory/%s"

	pwpCookieNameFmt = "PWP.ASPXAUTH=%s"

	pointHistoryJSONPath = "Data.1.Value"
	eventDataJSONPath    = "Data.Value"
)

func createInsecureClient() *http.Client {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	return client
}

func getPointHistory(dcinumber string) string {
	URL := fmt.Sprintf(pointHistoryURL, dcinumber)

	client := createInsecureClient()

	request, err := http.NewRequest("POST", URL, strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Fetching eventhistory for dcinumber %s\n", dcinumber)

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

	res := gjson.Get(string(body), pointHistoryJSONPath).String()

	return res
}

func fetchEventData(eventID string) string {
	client := createInsecureClient()

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
