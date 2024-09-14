package models

type Order struct {
    ID            string  `json:"id"`
    Symbol        string  `json:"symbol"`
    Quantity      int     `json:"quantity"`
    Price         float64 `json:"price"`
    Side          string  `json:"side"` // buy or sell
    Status        string  `json:"status"` // created, executed, canceled
    CreatedAt     int64   `json:"created_at"`
}
