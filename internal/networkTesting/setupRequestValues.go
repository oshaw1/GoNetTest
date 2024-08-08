package networkTesting

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/oshaw1/go-net-test/config"
)

func SetupHost(conf *config.Config, r *http.Request) (string, error) {
	host := conf.Host

	if queryHost := r.URL.Query().Get("host"); queryHost != "" {
		host = queryHost
	}

	if host == "" {
		return "", fmt.Errorf("host is required")
	}

	return host, nil
}

func SetupPort(conf *config.Config, r *http.Request) (int, error) {
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
