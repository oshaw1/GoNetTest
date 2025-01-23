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
	"github.com/oshaw1/go-net-test/internal/scheduler"
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
	ip := getOutboundIP()

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
	scheduler := scheduler.NewScheduler("http://" + ip + conf.Port)

	schedulerHandler := handler.NewSchedulerHandler(scheduler)
	networkTestHandler := handler.NewNetworkTestHandler(tester, repository)
	chartHandler := handler.NewChartHandler(repository, conf)
	utilHandler := &handler.UtilHandler{}
	dashboardHandler := handler.NewDashboardHandler(repository, "internal/pageGeneration/templates/*.tmpl")

	scheduler.Start()
	defer scheduler.Stop()
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

	mux.HandleFunc("/charts/generate", middleware.LoggingMiddleware(chartHandler.GenerateChart))
	mux.HandleFunc("/charts/generate-historic", middleware.LoggingMiddleware(chartHandler.GenerateHistoricChart))

	mux.HandleFunc("/schedule/create", middleware.LoggingMiddleware(schedulerHandler.HandleCreateSchedule))
	mux.HandleFunc("/schedule/list", middleware.LoggingMiddleware(schedulerHandler.HandleGetSchedules))

	server := &http.Server{
		Handler:      mux,
		Addr:         conf.Port,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	mux.PrintRoutes(ip, conf.Port)

	log.Printf("Server starting on http://%s%s\n", ip, conf.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
