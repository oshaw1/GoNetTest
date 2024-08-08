package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

func Load(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func SetupHost(conf *Config, r *http.Request) (string, error) {
	host := conf.Host

	if queryHost := r.URL.Query().Get("host"); queryHost != "" {
		host = queryHost
	}

	if host == "" {
		return "", fmt.Errorf("host is required")
	}

	return host, nil
}

func SetupPort(conf *Config, r *http.Request) (int, error) {
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
