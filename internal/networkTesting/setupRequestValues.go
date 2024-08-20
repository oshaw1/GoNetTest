package networkTesting

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/oshaw1/go-net-test/config"
)

func GetNetworkParams(r *http.Request) (string, int) {
	conf, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	host, err := setupHost(conf, r)
	if err != nil {
		log.Printf("Failed to load configuration host defaulting to localhost, Caused by: %v", err)
		return "localhost", 8080
	}

	port, err := setupPort(conf, r)
	if err != nil {
		log.Printf("Failed to load configuration port defaulting to localhost, Caused by: %v", err)
		return "localhost", 8080
	}
	return host, port
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

func setupPort(conf *config.Config, r *http.Request) (int, error) {
	port := conf.Port

	if queryPort := r.URL.Query().Get("port"); queryPort != "" {
		parsedPort, err := strconv.Atoi(queryPort)
		if err != nil {
			return port, err
		}
		port = parsedPort
	}

	return port, nil
}
