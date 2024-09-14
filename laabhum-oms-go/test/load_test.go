package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

type Order struct {
	Symbol   string  `json:"symbol"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Side     string  `json:"side"`
}
func TestLoadWithConcurrency(t *testing.T) {
	url := "http://localhost:8080/oms/scalper/order"
	concurrency := 10 // Reduced concurrency
	totalRequests := 1000

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrency; j++ {
				order := Order{
					Symbol:   "AAPL",
					Quantity: 100,
					Price:    150.0,
					Side:     "buy",
				}
				jsonOrder, err := json.Marshal(order)
				if err != nil {
					t.Errorf("Failed to marshal order: %v", err)
					continue
				}
				resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonOrder))
				if err != nil {
					t.Errorf("Request failed: %v", err)
					continue
				}
				resp.Body.Close()
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	fmt.Printf("Total time: %v\n", duration)
	fmt.Printf("Requests per second: %.2f\n", float64(totalRequests)/duration.Seconds())
}
