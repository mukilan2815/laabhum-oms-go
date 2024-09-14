package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
    "github.com/laabhum/laabhum-oms-go/api"
    "github.com/laabhum/laabhum-oms-go/repository"
    "github.com/laabhum/laabhum-oms-go/service"
)

func main() {
    repo := repository.NewOrderRepository()
    oms := service.NewOMS(repo)

    r := mux.NewRouter()
    api.RegisterHandlers(r, oms)

    r.Use(loggingMiddleware)

    log.Println("Server is running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Incoming request: Method=%s, URL=%s, RemoteAddr=%s", r.Method, r.URL.String(), r.RemoteAddr)

        rr := &responseRecorder{ResponseWriter: w}
        
        next.ServeHTTP(rr, r)

        log.Printf("Response status: %d", rr.statusCode)
    })
}

type responseRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
    rr.statusCode = statusCode
    rr.ResponseWriter.WriteHeader(statusCode)
}
