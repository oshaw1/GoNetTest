package networkTesting

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/config"
)

func GetHost(r *http.Request) string {
	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	host, err := setupHost(conf, r)
	if err != nil {
		log.Printf("Failed to load configuration host defaulting to localhost, Caused by: %v", err)
		return "localhost"
	}

	return host
}

func setupHost(conf *config.Config, r *http.Request) (string, error) {
	host := conf.Host

	if queryHost := r.URL.Query().Get("host"); queryHost != "" {
		host = queryHost
	}

	if host == "" {
		return "", fmt.Errorf("host is required")
	}

	return host, nil
}
