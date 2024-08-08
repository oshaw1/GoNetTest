package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/api/handler"
	"github.com/oshaw1/go-net-test/api/middleware"
)

func main() {
	// routes
	http.HandleFunc("/health", middleware.LoggingMiddleware(handler.HealthCheck))
	http.HandleFunc("/networktest", middleware.LoggingMiddleware(handler.NetworkTest))

	// server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
