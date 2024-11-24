package main

import (
	"fmt"
	"net/http"
)

type ClientRequestHandler struct{}

func (c *ClientRequestHandler) HandleRequest(r *http.Request) {
	// Handle incoming request from client
	// Here we just log the method and URL of the request for demonstration
	fmt.Println("Handling request:", r.Method, r.URL)
}
