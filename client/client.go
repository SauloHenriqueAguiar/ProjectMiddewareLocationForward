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
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main() {
	requestor := &Requestor{}
	payload := map[string]string{"key": "value"}
	response, err := requestor.SendRequest("http://localhost:8080/proxy", payload)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		defer response.Body.Close()
		fmt.Println("Response received:", response.Status)
	}
}
