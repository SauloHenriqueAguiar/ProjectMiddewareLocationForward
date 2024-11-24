package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Unmarshaler struct{}

func (u *Unmarshaler) Unmarshal(r *http.Request) (map[string]string, error) {
	var payload map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling request body: %v", err)
	}
	return payload, nil
}
