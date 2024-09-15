package unit

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/adapter"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/pkg/kafka"
)

func TestCreateOrder(t *testing.T) {
	producer := &kafka.MockProducer{}
	handler := adapter.NewOrderHandler(producer)

	body := []byte(`{"symbol": "AAPL", "quantity": 10, "price": 150.0}`)
	req, err := http.NewRequest("POST", "/order", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.CreateOrder(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}
