package main

import (
    "net/http"
    "sync"
    "testing"
    "time"
)

const numRequests = 10000

func performRequests(url string, wg *sync.WaitGroup, results chan<- time.Duration) {
    defer wg.Done()
    start := time.Now()
    resp, err := http.Get(url)
    if err != nil {
        results <- 0
        return
    }
    resp.Body.Close()
    duration := time.Since(start)
    results <- duration
}

func TestPerformance(t *testing.T) {
    gatewayURL := "http://localhost:8080/api/v1/oms/scalper/order"

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
    totalDuration := time.Since(start)

    var total time.Duration
    for result := range results {
        total += result
    }
    avgDuration := total / numRequests
    t.Logf("Total duration for %d requests: %v", numRequests, totalDuration)
    t.Logf("Average duration per request: %v", avgDuration)
}