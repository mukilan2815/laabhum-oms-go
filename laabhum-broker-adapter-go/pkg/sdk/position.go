package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func GetPositions(baseURL string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", baseURL+"/positions", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var positions []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, err
	}

	return positions, nil
}

func ConvertPosition(baseURL string, conversionPayload map[string]interface{}) error {
	payload, _ := json.Marshal(conversionPayload)
	req, err := http.NewRequest("POST", baseURL+"/position/convert", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	return err
}
