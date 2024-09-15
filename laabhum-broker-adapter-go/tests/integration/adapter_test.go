package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/adapter"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/kafka/consumer"
	"github.com/Mukilan-T/laabhum-broker-adapter-go/internal/kafka"
	"github.com/gorilla/mux"
)

func TestOrderFlow(t *testing.T) {
	producer := &kafka.MockProducer{}
	router := mux.NewRouter()
	consumer := &kafka.MockConsumer{} // Ensure MockConsumer is defined in the kafka package
	adapter.SetupRoutes(router, producer, consumer)

	req, _ := http.NewRequest("POST", "/order", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}
