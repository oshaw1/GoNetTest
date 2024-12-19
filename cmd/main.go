package main

import (
	"log"
	"net/http"
	"os"

	"github.com/oshaw1/go-net-test/api/handler"
	"github.com/oshaw1/go-net-test/api/middleware"
	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func initDataDir() {
	err := os.MkdirAll("data/output/", 0755)
	if err != nil {
		log.Fatalf("failed to create directory structure: %e", err)
	}
}

func main() {
	initDataDir()

	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	repository := dataManagement.NewRepository("data/output", conf)
	tester := networkTesting.NewNetworkTester(conf)

	networkTestHandler := handler.NewNetworkTestHandler(tester, repository)
	utilHandler := &handler.UtilHandler{}
	dashboardHandler := handler.NewDashboardHandler(repository, "internal/pageGeneration/templates/*.tmpl")

	fs := http.FileServer(http.Dir("web/static"))
	// file server
	http.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data/output"))))
	http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))

	http.HandleFunc("/health", middleware.LoggingMiddleware(utilHandler.HealthCheck))

	http.HandleFunc("/dashboard/", middleware.LoggingMiddleware(dashboardHandler.ServeDashboard))
	http.HandleFunc("/dashboard/recent-tests-quadrant", middleware.LoggingMiddleware(dashboardHandler.GetRecentQuadrant))
	http.HandleFunc("/dashboard/chart", middleware.LoggingMiddleware(dashboardHandler.GetChart))
	http.HandleFunc("/dashboard/data", middleware.LoggingMiddleware(dashboardHandler.GetData))

	http.HandleFunc("/networktest", middleware.LoggingMiddleware(networkTestHandler.HandleNetworkTest))
	http.HandleFunc("/networktest/test-results", middleware.LoggingMiddleware(networkTestHandler.GetResults))

	// server
	port := ":7000"
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
