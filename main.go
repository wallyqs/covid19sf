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

type SanFranciscoData struct {
	// TestPositivityRate
	TestPositivityRate string `json:"pct"`

	// Positive is the number of cases which are positive.
	Positive string `json:"pos"`

	// ResultData
	ResultDate string `json:"result_date"`

	// Tests
	Tests string `json:"tests"`
}

func

func main() {
	resp, err := http.Get(dataSFURL)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	var recentData []*SanFranciscoData
	err = json.Unmarshal(data, &recentData)
	if err != nil {
		log.Fatal(err)
	}

	if len(recentData) < 1 {
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

	fmt.Printf("|----------------------|-----------------|-----------------|----------------------|\n")
	fmt.Printf("| City                 | Positive  Cases | Negative Cases  | Test Positivity Rate |\n")
	fmt.Printf("|----------------------|-----------------|-----------------|----------------------|\n")
	fmt.Printf("| San Francisco        | %-15d | %-15d | %-20.4f |\n",
		pos,
		neg,
		tpr,
	)
	fmt.Printf("|----------------------|-----------------|-----------------|----------------------|\n")
}
