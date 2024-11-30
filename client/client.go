package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// Marshaller é responsável por serializar o payload.
type Marshaller struct{}

func (m *Marshaller) Marshal(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

// Requestor é responsável por enviar requisições HTTP e lidar com redirecionamentos.
type Requestor struct {
	currentURL string
}

func NewRequestor(initialURL string) *Requestor {
	return &Requestor{currentURL: initialURL}
}

func (r *Requestor) SendRequest(payload interface{}) (*http.Response, error) {
	marshaller := &Marshaller{}
	data, err := marshaller.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data: %v", err)
	}

	resp, err := http.Post(r.currentURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}

	// Se a resposta for 301 ou 302, atualiza a URL e reenvia a requisição.
	if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound {
		newURL := resp.Header.Get("Location")
		if newURL == "" {
			return nil, fmt.Errorf("redirection without Location header")
		}
		fmt.Printf("Redirected to new URL: %s\n", newURL)
		r.currentURL = newURL
		resp.Body.Close() // Fecha a resposta anterior antes de reenviar
		return r.SendRequest(payload)
	}

	return resp, nil
}

// Invoker encapsula a lógica de requisições, usando o Requestor internamente.
type Invoker struct {
	requestor *Requestor
}

func (i *Invoker) InvokeRequest(payload interface{}) (*http.Response, error) {
	return i.requestor.SendRequest(payload)
}

// Função para gerar dados aleatórios
func generateRandomPayload() map[string]interface{} {
	randomValue := rand.Intn(100)
	return map[string]interface{}{
		"key": randomValue,
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Inicializa o Requestor com a URL inicial do proxy
	initialURL := "http://localhost:8082/proxy"
	requestor := NewRequestor(initialURL)
	invoker := &Invoker{requestor: requestor}

	for {
		// Gera dados aleatórios para enviar
		payload := generateRandomPayload()

		// Envia a requisição usando o Invoker
		response, err := invoker.InvokeRequest(payload)
		if err != nil {
			fmt.Println("Error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Lê a resposta do servidor
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Response received: %s\n", string(body))

		// Aguarda 3 segundos antes de enviar a próxima requisição
		time.Sleep(3 * time.Second)
	}
}
