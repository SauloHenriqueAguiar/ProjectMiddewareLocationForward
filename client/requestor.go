package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Requestor struct{}

func (r *Requestor) SendRequest(url string, payload interface{}) (*http.Response, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	return resp, nil
}
