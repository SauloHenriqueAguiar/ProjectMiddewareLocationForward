package main

import (
	"fmt"
	"net/http"
)

type ServerRequestHandler struct{}

func (s *ServerRequestHandler) HandleRequest(r *http.Request) {
	// Handle incoming request (for example, log the request details)
	fmt.Println("Received request:", r.Method, r.URL)
}
