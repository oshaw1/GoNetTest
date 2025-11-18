package main

import (
	"fmt"
	"log"
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
=================================================
     GoNetTest - Network Testing Made Simple
     Created by Owen Shaw
     GitHub: github.com/oshaw1
=================================================
    `
	fmt.Println(banner)
}

func main() {
	printBanner()
	ip := "0.0.0.0"

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
	scheduler := scheduler.NewScheduler("http://"+ip+conf.Port, conf.Scheduler.Schedule)

	schedulerHandler := handler.NewSchedulerHandler(scheduler)
	networkTestHandler := handler.NewNetworkTestHandler(tester, repository)
	chartHandler := handler.NewChartHandler(repository, conf)
	utilHandler := &handler.UtilHandler{}
	dashboardHandler := handler.NewDashboardHandler(repository, "internal/pageGeneration/templates/*.gohtml", scheduler)

	mux := middleware.NewRouteMux()

	mux.Handle("/data/output/", http.StripPrefix("/data/output/", http.FileServer(http.Dir("data/output"))))
	mux.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("/health", middleware.LoggingMiddleware(utilHandler.HealthCheck))

	mux.HandleFunc("/dashboard/", middleware.LoggingMiddleware(dashboardHandler.ServeDashboard))
	mux.HandleFunc("/dashboard/tests", middleware.LoggingMiddleware(dashboardHandler.ServeTestQuadrant))
	mux.HandleFunc("/dashboard/schedule", middleware.LoggingMiddleware(dashboardHandler.ServeSchedule))

	mux.HandleFunc("/networktest", middleware.LoggingMiddleware(networkTestHandler.HandleNetworkTest))
	mux.HandleFunc("/networktest/test-results", middleware.LoggingMiddleware(networkTestHandler.GetResults))

	mux.HandleFunc("/charts/generate", middleware.LoggingMiddleware(chartHandler.GenerateChart))
	mux.HandleFunc("/charts/generate-historic", middleware.LoggingMiddleware(chartHandler.GenerateHistoricChart))

	mux.HandleFunc("/schedule/create", middleware.LoggingMiddleware(schedulerHandler.HandleCreateSchedule))
	mux.HandleFunc("/schedule/list", middleware.LoggingMiddleware(schedulerHandler.HandleGetSchedule))
	mux.HandleFunc("/schedule/export", middleware.LoggingMiddleware(schedulerHandler.HandleExportSchedule))
	mux.HandleFunc("/schedule/import", middleware.LoggingMiddleware(schedulerHandler.HandleImportSchedule))
	mux.HandleFunc("/schedule/delete", middleware.LoggingMiddleware(schedulerHandler.HandleDeleteSchedule))
	mux.HandleFunc("/schedule/edit", middleware.LoggingMiddleware(schedulerHandler.HandleEditSchedule))

	server := &http.Server{
		Handler:      mux,
		Addr:         conf.Port,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
		IdleTimeout:  120 * time.Second,
	}

	mux.PrintRoutes(ip, conf.Port)

	scheduler.Start()
	defer scheduler.Stop()
	log.Printf("Server accessible on http://%s%s\n", ip, conf.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
