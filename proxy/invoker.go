package main

import (
	"bytes"
	"fmt"
	"net/http"
)

type Invoker struct{}

func (i *Invoker) InvokeRequest(url string, payload interface{}) (*http.Response, error) {
	marshaller := &Marshaller{}
	data, err := marshaller.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error invoking request: %v", err)
	}

	return resp, nil
}
