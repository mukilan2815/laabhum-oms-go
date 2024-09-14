package models

type Order struct {
	ID        string  `json:"id"`
	Symbol    string  `json:"symbol"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Side      string  `json:"side"` // "buy" or "sell"
	Status    string  `json:"status"`
	CreatedAt int64   `json:"created_at"`
}

type ScalperOrder struct {
	ID          string  `json:"id"`
	ParentOrder Order   `json:"parent_order"`
	ChildOrders []Order `json:"child_orders"`
	Status      string  `json:"status"`
	CreatedAt   int64   `json:"created_at"`
}

type Trade struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}