package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/oshaw1/go-net-test/api/handler"
	"github.com/oshaw1/go-net-test/api/middleware"
	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/dataManagement"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
)

func printBanner() {
	banner := `
    =============================================
     goNetTest - Network Testing Made Simple
     Created by Owen Shaw
     GitHub: github.com/oshaw1
    =============================================
    `
	fmt.Println(banner)
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func main() {
	printBanner()

	err := os.MkdirAll("data/output/", 0755)
	if err != nil {
		log.Fatalf("failed to create directory structure: %e", err)
	}

	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	repository := dataManagement.NewRepository("data/output", conf)
	tester := networkTesting.NewNetworkTester(conf)

	networkTestHandler := handler.NewNetworkTestHandler(tester, repository)
	utilHandler := &handler.UtilHandler{}
	dashboardHandler := handler.NewDashboardHandler(repository, "internal/pageGeneration/templates/*.tmpl")

	mux := middleware.NewRouteMux()

	fs := http.FileServer(http.Dir("web/static"))

	mux.Handle("/data/", http.StripPrefix("/data/", http.FileServer(http.Dir("data/output"))))
	mux.Handle("/web/static/", http.StripPrefix("/web/static/", fs))

	mux.HandleFunc("/health", middleware.LoggingMiddleware(utilHandler.HealthCheck))

	mux.HandleFunc("/dashboard/", middleware.LoggingMiddleware(dashboardHandler.ServeDashboard))
	mux.HandleFunc("/dashboard/recent-tests-quadrant", middleware.LoggingMiddleware(dashboardHandler.GetRecentQuadrant))
	mux.HandleFunc("/dashboard/chart", middleware.LoggingMiddleware(dashboardHandler.GetChart))
	mux.HandleFunc("/dashboard/data", middleware.LoggingMiddleware(dashboardHandler.GetData))

	mux.HandleFunc("/networktest", middleware.LoggingMiddleware(networkTestHandler.HandleNetworkTest))
	mux.HandleFunc("/networktest/test-results", middleware.LoggingMiddleware(networkTestHandler.GetResults))

	port := ":7000"
	server := &http.Server{
		Handler:      mux,
		Addr:         port,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	ip := getOutboundIP()
	mux.PrintRoutes(ip, port)

	log.Printf("Server starting on http://%s%s\n", ip, port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
