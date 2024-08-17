package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/config"
	"github.com/oshaw1/go-net-test/internal/networkTesting"
	pagegeneration "github.com/oshaw1/go-net-test/internal/pageGeneration"
)

type PageHandler struct {
}

func (h PageHandler) HandleICMPResult(w http.ResponseWriter, r *http.Request) {
	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	host, err := networkTesting.SetupHost(conf, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	port, err := networkTesting.SetupPort(conf, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := networkTesting.TestNetwork(host, port)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error performing network test: %v", err), http.StatusInternalServerError)
		return
	}

	html, err := pagegeneration.GenerateICMBTestResultHTML(result)
	if err != nil {
		http.Error(w, "Error generating result", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
