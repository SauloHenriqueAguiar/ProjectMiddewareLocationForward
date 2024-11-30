package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Marshaller é responsável por serializar e desserializar payloads.
type Marshaller struct{}

func (m *Marshaller) Marshal(payload interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (m *Marshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Requestor é responsável por realizar requisições HTTP externas.
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

// Invoker é responsável por encapsular a lógica de requisições e respostas.
type Invoker struct {
	requestor *Requestor
}

func (i *Invoker) InvokeRequest(url string, payload interface{}) ([]byte, error) {
	resp, err := i.requestor.SendRequest(url, payload)
	if err != nil {
		return nil, fmt.Errorf("error invoking request: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return data, nil
}

// ClientRequestHandler é responsável por lidar com requisições recebidas do cliente.
type ClientRequestHandler struct {
	invoker *Invoker
}

func (h *ClientRequestHandler) HandleClientRequest(w http.ResponseWriter, r *http.Request) {
	marshaller := &Marshaller{}
	var payload map[string]interface{}

	// Lê o corpo da requisição do cliente.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
		return
	}

	// Exibe a mensagem que foi recebida no proxy
	log.Printf("Proxy received request: %s", body)

	// Desserializa o payload recebido.
	err = marshaller.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error unmarshalling request body: %v", err), http.StatusBadRequest)
		return
	}

	// Faz a requisição para o server local.
	responseData, err := h.invoker.InvokeRequest("http://localhost:8081/process", payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error invoking server request: %v", err), http.StatusInternalServerError)
		return
	}

	// Exibe a resposta recebida do server (após a requisição)
	log.Printf("Proxy received response: %s", responseData)

	// Retorna a resposta para o cliente.
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func main() {
	// Inicializa os componentes.
	requestor := &Requestor{}
	invoker := &Invoker{requestor: requestor}
	handler := &ClientRequestHandler{invoker: invoker}

	// Configura o servidor HTTP.
	http.HandleFunc("/proxy", handler.HandleClientRequest)
	fmt.Println("Proxy server is running on port 8082...")
	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
