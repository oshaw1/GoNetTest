package main

import (
	"log"
	"net/http"
	"os"

	"github.com/oshaw1/go-net-test/api/handler"
	"github.com/oshaw1/go-net-test/api/middleware"
)

func initDataDir() {
	err := os.MkdirAll("data/output/", 0755)
	if err != nil {
		log.Fatalf("failed to create directory structure: %e", err)
	}
}

func main() {
	initDataDir()
	fs := http.FileServer(http.Dir("web/static"))
	// file server
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data/output"))))
	http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))

	networkTestHandler := &handler.NetworkTestHandler{}
	utilHandler := &handler.UtilHandler{}
	pageHandler := &handler.PageHandler{}

	http.HandleFunc("/health", middleware.LoggingMiddleware(utilHandler.HealthCheck))

	http.HandleFunc("/dashboard/", middleware.LoggingMiddleware(pageHandler.ServeDashboard))
	http.HandleFunc("/dashboard/runtest/icmb", middleware.LoggingMiddleware(pageHandler.RunICMBTest))
	http.HandleFunc("/dashboard/recent-tests-quadrant", middleware.LoggingMiddleware(pageHandler.GetRecentQuadrant))

	http.HandleFunc("/networktest/icmp", middleware.LoggingMiddleware(networkTestHandler.HandleICMPNetworkTest))

	// server
	port := ":7000"
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
