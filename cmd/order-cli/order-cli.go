package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strings" // Missing import for strings
	"sync"
	"time"
)

type Order struct {
    Symbol   string `json:"symbol"`
    Quantity int    `json:"quantity"`
}

func sendOrder(order Order, wg *sync.WaitGroup) {
    defer wg.Done()

    jsonData, err := json.Marshal(order)
    if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        return
    }

    resp, err := http.Post("http://localhost:8080/orders", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error sending request:", err)
        return
    }
    defer resp.Body.Close()

    fmt.Println("Order placed successfully. Status:", resp.Status)
}

func generateRandomOrder(symbols []string) Order {
    // Generate random quantity between 1 and 100
    quantity := rand.Intn(1000) + 1
    // Pick a random symbol
    symbol := symbols[rand.Intn(len(symbols))]
    return Order{
        Symbol:   symbol,
        Quantity: quantity,
    }
}

func main() {
    symbolsFlag := flag.String("symbols", "", "Comma-separated list of stock symbols")
    concurrency := flag.Int("concurrency", 1, "Number of concurrent requests")
    totalRequests := flag.Int("requests", 1, "Total number of requests to send")
    flag.Parse()

    if *symbolsFlag == "" || *concurrency <= 0 || *totalRequests <= 0 {
        fmt.Println("Usage: order-cli -symbols <symbols> -concurrency <number> -requests <number>")
        return
    }

    symbols := splitSymbols(*symbolsFlag)

    // Seed random number generator
    rand.Seed(time.Now().UnixNano())

    var wg sync.WaitGroup
    for i := 0; i < *totalRequests; i++ {
        wg.Add(1)
        go func() {
            order := generateRandomOrder(symbols)
            sendOrder(order, &wg)
        }()
        // Control the rate of sending requests
        time.Sleep(10 * time.Millisecond) // Adjust as needed
    }

    wg.Wait()
}

func splitSymbols(symbols string) []string {
    return strings.Split(symbols, ",")
}
