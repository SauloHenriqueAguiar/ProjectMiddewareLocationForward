package main

import "encoding/json"

type Marshaller struct{}

func (m *Marshaller) Marshal(payload interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}
