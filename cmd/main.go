package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

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
	http.HandleFunc("/dashboard/", middleware.LoggingMiddleware(serveHTML))
	http.HandleFunc("/health", middleware.LoggingMiddleware(utilHandler.HealthCheck))
	http.HandleFunc("/networktest/icmp", middleware.LoggingMiddleware(networkTestHandler.HandleICMPNetworkTest))
	http.HandleFunc("/page/icmpreturn", middleware.LoggingMiddleware(pageHandler.HandleICMPResult))

	// server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dashboard/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(filepath.Join("web", "static", "dashboard.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
