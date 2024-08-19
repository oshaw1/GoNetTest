package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/api/handler"
	"github.com/oshaw1/go-net-test/api/middleware"
)

func main() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))

	networkTestHandler := &handler.NetworkTestHandler{}
	utilHandler := &handler.UtilHandler{}
	pageHandler := &handler.PageHandler{}

	// routes
	http.HandleFunc("/health", middleware.LoggingMiddleware(utilHandler.HealthCheck))

	http.HandleFunc("/dashboard/", middleware.LoggingMiddleware(pageHandler.ServeDashboard))
	http.HandleFunc("/dashboard/recent-tests-data", middleware.LoggingMiddleware(pageHandler.GetRecentTestData))
	http.HandleFunc("/dashboard/runtest/icmb", middleware.LoggingMiddleware(pageHandler.RunICMBTest))
	//http.HandleFunc("/dashboard/recent-tests-charts", middleware.LoggingMiddleware(pageHandler.GetRecentTestCharts))

	http.HandleFunc("/networktest/icmp", middleware.LoggingMiddleware(networkTestHandler.HandleICMPNetworkTest))

	// server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
