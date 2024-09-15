package broker

type OrderRequest struct {
	Symbol   string  `json:"symbol"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type PositionResponse struct {
	Symbol string  `json:"symbol"`
	Qty    int     `json:"qty"`
	Price  float64 `json:"price"`
}
