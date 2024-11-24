package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	ServerRequestHandler *ServerRequestHandler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Process the incoming request using ServerRequestHandler
	serverRequestHandler := &ServerRequestHandler{}
	serverRequestHandler.HandleRequest(r)

	// Unmarshal the request payload to a map
	unmarshaler := &Unmarshaler{}
	payload, err := unmarshaler.Unmarshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process the payload and create a response
	response := map[string]string{"received": "true", "data": payload["key"]}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	server := &Server{}
	http.Handle("/server", server)
	fmt.Println("Server running on port 8082")
	http.ListenAndServe(":8082", nil)
}
