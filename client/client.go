package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// Marshaller é responsável por serializar o payload.
type Marshaller struct{}

func (m *Marshaller) Marshal(payload interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Requestor é responsável por enviar requisições HTTP.
type Requestor struct{}

func (r *Requestor) SendRequest(url string, payload interface{}) (*http.Response, error) {
	marshaller := &Marshaller{}
	data, err := marshaller.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	return resp, nil
}

// Invoker encapsula a lógica de requisições, usando o Requestor internamente.
type Invoker struct {
	requestor *Requestor
}

func (i *Invoker) InvokeRequest(url string, payload interface{}) (*http.Response, error) {
	return i.requestor.SendRequest(url, payload)
}

// Função para gerar dados aleatórios
func generateRandomPayload() map[string]interface{} {
	// Gerando um valor aleatório para "key"
	randomValue := rand.Intn(100) // Número aleatório entre 0 e 99
	return map[string]interface{}{
		"key": randomValue,
	}
}

func main() {
	// Inicializa os componentes
	requestor := &Requestor{}
	invoker := &Invoker{requestor: requestor}

	for {
		// Gera dados aleatórios
		payload := generateRandomPayload()

		// Envia a requisição usando o Invoker
		response, err := invoker.InvokeRequest("http://localhost:8082/proxy", payload)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Exibe a resposta
		defer response.Body.Close()
		fmt.Println("Response received:", response.Status)

		// Aguarda 3 segundos antes de enviar a próxima requisição
		time.Sleep(3 * time.Second)
	}
}
