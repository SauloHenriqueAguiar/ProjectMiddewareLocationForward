package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Proxy struct {
	Requestor            *Requestor
	ClientRequestHandler *ClientRequestHandler
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Process the incoming request using ClientRequestHandler
	clientRequestHandler := &ClientRequestHandler{}
	clientRequestHandler.HandleRequest(r)

	// Get the payload (here we mock the data, but this would be a real request processing)
	payload := map[string]string{"key": "value"}

	// Use the Requestor to forward the request to the server
	requestor := &Requestor{}
	response, err := requestor.SendRequest("http://localhost:8082/server", payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Return the response from the server
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(body)
}

func main() {
	proxy := &Proxy{}
	http.Handle("/proxy", proxy)
	fmt.Println("Proxy server running on port 8081")
	http.ListenAndServe(":8081", nil)
}
