package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

const (
	// dataSFURL is the endpoint from where to get the TPR data
	// from SanFrancisco.
	dataSFURL = "https://data.sfgov.org/resource/nfpa-mg4g.json"

	// dataSFCasesURL is the endpoint from where can get the stats
	// about the cases in San Francisco.
	dataSFCasesURL = "https://data.sfgov.org/resource/tvq9-ec9w.json"
)

const (
	version     = "0.1.0"
	releaseDate = "April 10th, 2020"
)

type SanFranciscoTPRData struct {
	// TestPositivityRate
	TestPositivityRate string `json:"pct"`

	// Positive is the number of cases which are positive.
	Positive string `json:"pos"`

	// ResultData
	ResultDate string `json:"result_date"`

	// Tests
	Tests string `json:"tests"`
}

type SanFranciscoCasesData struct {
	Date                 string `json:"date"`
	TransmissionCategory string `json:"transmission_category"`
	CaseDisposition      string `json:"case_disposition"`
	CaseCount            string `json:"case_count"`
}

func curl(endpoint string) ([]byte, error) {
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func main() {
	data, err := curl(dataSFURL)
	if err != nil {
		log.Fatal(err)
	}

	var recentData []*SanFranciscoTPRData
	err = json.Unmarshal(data, &recentData)
	if err != nil {
		log.Fatal(err)
	}
	if len(recentData) < 1 {
		log.Fatal("No data!")
	}

	var casesData []*SanFranciscoCasesData
	data, err = curl(dataSFCasesURL)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &casesData)
	if err != nil {
		log.Fatal(err)
	}
	if len(casesData) < 1 {
		log.Fatal("No data!")
	}

	// Get the most recent data from datasf.
	latest := recentData[0]
	pos, err := strconv.Atoi(latest.Positive)
	if err != nil {
		log.Fatal(err)
	}
	tests, err := strconv.Atoi(latest.Tests)
	if err != nil {
		log.Fatal(err)
	}
	tpr, err := strconv.ParseFloat(latest.TestPositivityRate, 64)
	if err != nil {
		log.Fatal(err)
	}
	neg := tests - pos

	// Cases data
	totalDeaths := 0
	totalPositive := 0
	for _, cs := range casesData {
		if cs.CaseDisposition == "Death" {
			n, err := strconv.Atoi(cs.CaseCount)
			if err != nil {
				log.Fatal(err)
			}
			totalDeaths += n
		} else if cs.CaseDisposition == "Confirmed" {
			n, err := strconv.Atoi(cs.CaseCount)
			if err != nil {
				log.Fatal(err)
			}
			totalPositive += n
		}
	}
	fmt.Printf("|----------------------|-----------------------|-----------------------|----------------------|----------------|----------------|\n")
	fmt.Printf("| City                 | Positive (Past 24hrs) | Negative (Past 24hrs) | Test Positivity Rate | Total Positive | Total Deaths   |\n")
	fmt.Printf("|----------------------|-----------------------|-----------------------|----------------------|----------------|----------------|\n")
	fmt.Printf("| San Francisco        | %-21d | %-21d | %-20.4f | %-14d | %-14d |\n",
		pos,
		neg,
		tpr,
		totalPositive,
		totalDeaths,
	)
	fmt.Printf("|----------------------|-----------------------|-----------------------|----------------------|----------------|----------------|\n")
}
