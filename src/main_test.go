package main

import (
	"encoding/json"
	"fmt"
	"git.garena.com/shopee/insurance/insurance-backend/insurance-hub/moving-window/src/api"
	"git.garena.com/shopee/insurance/insurance-backend/insurance-hub/moving-window/src/common"
	"git.garena.com/shopee/insurance/insurance-backend/insurance-hub/moving-window/src/util"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	concurrencyThread   = 50
	bulkPerThread       = 20
	testRoundsInSeconds = 70 // Should be greater than 60s to cover the edge case
)

func TestConcurrencyCounts(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	totalBaseCount := makeRequest(t, "http://localhost:8081/api/requests")
	time.Sleep(time.Second)

	for i := 1; i <= testRoundsInSeconds; i++ {
		roundIndex := i

		<-ticker.C
		var wg sync.WaitGroup
		for j := 0; j < concurrencyThread; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for k := 0; k < bulkPerThread; k++ {
					makeRequest(t, "http://localhost:8081/api/requests")
				}
			}()
		}
		wg.Wait()
		currentCnt := makeRequest(t, "http://localhost:8081/api/requests")

		expectedCnt := totalBaseCount + (concurrencyThread*bulkPerThread+1)*(int64(roundIndex))
		if roundIndex < (common.WindowSize - 1) {
			if currentCnt != expectedCnt {
				t.Errorf("Round: %v, currentCnt %v, expectedCnt %v", roundIndex, currentCnt, expectedCnt)
				return
			}
		} else {
			// This is a not accurate number, because the server is recollecting the counts buckets from the oldest round.
			expectedCnt = (concurrencyThread*bulkPerThread + 1) * (common.WindowSize - 1)
		}
		fmt.Printf("Round: %v, totalBaseCount %v, currentCnt %v, expectedCnt %v \n", roundIndex, totalBaseCount, currentCnt, expectedCnt)
	}

	fmt.Println("Test passed")
}

func TestAfterRebootServer(t *testing.T) {
	// Need reboot server first
	currentCnt := makeRequest(t, "http://localhost:8081/api/requests")
	fmt.Printf("Reboot server: %v\n", currentCnt)

}

func makeRequest(t *testing.T, url string) int64 {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	util.CheckErr(err)

	var response Response
	json.Unmarshal(body, &response)
	util.CheckErr(err)
	return response.Data.TotalRequest
}

type Response struct {
	Status string        `json:"status"`
	Code   int           `json:"code"`
	Data   api.CountDada `json:"data"`
}
