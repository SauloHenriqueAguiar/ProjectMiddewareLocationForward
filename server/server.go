package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Unmarshaller é responsável por desserializar dados.
type Unmarshaller struct{}

func (u *Unmarshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// ServerRequestHandler lida com as requisições recebidas.
type ServerRequestHandler struct {
	unmarshaller    *Unmarshaller
	redirectEnabled bool
	newAddress      string
}

// HandleRequest processa a requisição ou redireciona se configurado.
func (h *ServerRequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if h.redirectEnabled {
		// Retorna o redirecionamento com o novo endereço
		log.Printf("Redirecting request to new address: %s", h.newAddress)
		w.Header().Set("Location", h.newAddress)
		http.Error(w, "Server moved", http.StatusMovedPermanently)
		return
	}

	// Lê o corpo da requisição
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Server received request: %s", body)

	var payload map[string]interface{}
	err = h.unmarshaller.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error unmarshalling request body: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Server processed request, payload: %+v", payload)

	// Cria e retorna a resposta
	response := map[string]interface{}{
		"status":  "success",
		"payload": payload,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func main() {
	unmarshaller := &Unmarshaller{}
	redirectEnabled := false
	newAddress := os.Getenv("NEW_SERVER_ADDRESS") // Defina o novo endereço via variável de ambiente

	if newAddress != "" {
		redirectEnabled = true
	}

	handler := &ServerRequestHandler{
		unmarshaller:    unmarshaller,
		redirectEnabled: redirectEnabled,
		newAddress:      newAddress,
	}

	http.HandleFunc("/process", handler.HandleRequest)
	port := "8081"
	log.Printf("Server is running on port %s...", port)

	go func() {
		if redirectEnabled {
			// Simula migração após 30 segundos (para testes)
			time.Sleep(30 * time.Second)
			handler.redirectEnabled = true
			handler.newAddress = "http://new-server:8081/process"
			log.Println("Server is now redirecting requests.")
		}
	}()

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
