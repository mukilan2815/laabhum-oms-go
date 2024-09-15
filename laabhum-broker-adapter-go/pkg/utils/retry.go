package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	for i := 0; i < attempts; i++ {
		if err := fn(); err != nil {
			time.Sleep(sleep)
			sleep = sleep * 2
			continue
		}
		return nil
	}
	return errors.New("all attempts failed")
}
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
    RespondWithJSON(w, code, map[string]string{"error": message})
}
