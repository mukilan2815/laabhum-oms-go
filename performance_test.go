package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

const (
	numRequests = 1000
	concurrency = 50
)

func performRequests(url string, wg *sync.WaitGroup, results chan<- time.Duration) {
	defer wg.Done()

	start := time.Now()
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Printf("Failed to perform request: %v\n", err)
		return
	}
	resp.Body.Close()
	duration := time.Since(start)
	results <- duration
}

func TestPerformance(t *testing.T) {
	gatewayURL := "http://localhost:8080/api/v1/oms/scalper/order"
	omsURL := "http://localhost:8081/api/v1/oms/order"

	var wg sync.WaitGroup
	results := make(chan time.Duration, numRequests)

	// Test Gateway
	wg.Add(numRequests)
	start := time.Now()
	for i := 0; i < numRequests; i++ {
		go performRequests(gatewayURL, &wg, results)
	}
	wg.Wait()
	close(results)
	gatewayDurations := make([]time.Duration, 0, numRequests)
	for duration := range results {
		gatewayDurations = append(gatewayDurations, duration)
	}
	gatewayTotalTime := time.Since(start)
	gatewayAverageTime := time.Duration(float64(gatewayTotalTime) / float64(numRequests))

	// Test OMS
	results = make(chan time.Duration, numRequests)
	wg.Add(numRequests)
	start = time.Now()
	for i := 0; i < numRequests; i++ {
		go performRequests(omsURL, &wg, results)
	}
	wg.Wait()
	close(results)
	omsDurations := make([]time.Duration, 0, numRequests)
	for duration := range results {
		omsDurations = append(omsDurations, duration)
	}
	omsTotalTime := time.Since(start)
	omsAverageTime := time.Duration(float64(omsTotalTime) / float64(numRequests))

	// Calculate metrics
	fmt.Printf("Gateway took %v for %d requests\n", gatewayTotalTime, numRequests)
	fmt.Printf("Gateway average response time: %v\n", gatewayAverageTime)
	fmt.Printf("OMS took %v for %d requests\n", omsTotalTime, numRequests)
	fmt.Printf("OMS average response time: %v\n", omsAverageTime)
}
