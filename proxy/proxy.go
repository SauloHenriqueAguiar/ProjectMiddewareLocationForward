package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// Marshaller é responsável por serializar e desserializar payloads.
type Marshaller struct{}

func (m *Marshaller) Marshal(payload interface{}) ([]byte, error) {
	return json.Marshal(payload)
}

func (m *Marshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Requestor é responsável por enviar requisições HTTP externas.
type Requestor struct {
	serverAddress string
	mu            sync.Mutex
}

func NewRequestor(initialServerAddress string) *Requestor {
	return &Requestor{serverAddress: initialServerAddress}
}

// Atualiza a localização do servidor remoto
func (r *Requestor) UpdateServerAddress(newAddress string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.serverAddress = newAddress
	log.Printf("Server address updated to: %s", newAddress)
}

// Retorna a localização atual do servidor remoto
func (r *Requestor) GetServerAddress() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.serverAddress
}

// Envia requisição ao servidor remoto
func (r *Requestor) SendRequest(payload interface{}) (*http.Response, error) {
	marshaller := &Marshaller{}
	data, err := marshaller.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request data: %v", err)
	}

	url := fmt.Sprintf("http://%s/process", r.GetServerAddress())
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error sending request to %s: %v", url, err)
	}

	// Checa se a resposta indica redirecionamento
	if resp.StatusCode == http.StatusMovedPermanently || resp.StatusCode == http.StatusFound {
		newLocation := resp.Header.Get("Location")
		if newLocation == "" {
			return nil, fmt.Errorf("redirection without Location header")
		}
		log.Printf("Redirect received from server: %s", newLocation)
		r.UpdateServerAddress(newLocation)
		resp.Body.Close()
		return r.SendRequest(payload)
	}

	return resp, nil
}

// ClientRequestHandler lida com requisições do cliente.
type ClientRequestHandler struct {
	requestor *Requestor
}

func (h *ClientRequestHandler) HandleClientRequest(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
		return
	}

	marshaller := &Marshaller{}
	err = marshaller.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error unmarshalling request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Proxy received request: %s", body)

	// Envia a requisição ao servidor remoto através do Requestor
	response, err := h.requestor.SendRequest(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error forwarding request: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Copia a resposta para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
}

func main() {
	// Endereço inicial do servidor remoto
	initialServer := "localhost:8081" // Substituir pelo endereço IP do container em ambiente real
	requestor := NewRequestor(initialServer)
	handler := &ClientRequestHandler{requestor: requestor}

	// Configura o servidor HTTP do Proxy
	http.HandleFunc("/proxy", handler.HandleClientRequest)
	log.Println("Proxy server is running on port 8082...")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatalf("Error starting proxy server: %v", err)
	}
}
