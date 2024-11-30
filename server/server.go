package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Unmarshaller struct{}

// Método para deserializar os dados
func (u *Unmarshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type ServerRequestHandler struct {
	unmarshaller *Unmarshaller
}

// Método para lidar com as requisições recebidas
func (h *ServerRequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {

	// Lê o corpo da requisição
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading request body: %v", err), http.StatusBadRequest)
		return
	}

	// Exibe a mensagem recebida no servidor
	log.Printf("Server received request: %s", body)

	// Desserializa o payload recebido
	var payload map[string]interface{}
	err = h.unmarshaller.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error unmarshalling request body: %v", err), http.StatusBadRequest)
		return
	}

	// Verificação de que o payload foi recebido e processado corretamente
	log.Printf("Server processed request, payload: %+v", payload)

	// Cria a resposta
	response := map[string]interface{}{
		"status":  "success",
		"payload": payload,
	}

	// Retorna a resposta para o cliente
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func main() {

	unmarshaller := &Unmarshaller{}
	handler := &ServerRequestHandler{unmarshaller: unmarshaller}

	// Configura o servidor HTTP
	http.HandleFunc("/process", handler.HandleRequest)
	fmt.Println("Server is running on port 8081...")

	// Inicia o servidor
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
