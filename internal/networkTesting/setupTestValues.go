package networkTesting

import (
	"fmt"
	"log"
	"net/http"

	"github.com/oshaw1/go-net-test/config"
)

func GetHost(r *http.Request) string {
	conf := getConfig()

	host, err := setupHost(conf, r)
	if err != nil {
		log.Printf("Failed to load configuration host defaulting to localhost, Caused by: %v", err)
		return "localhost"
	}

	return host
}

func getTestValue[T any](valueName string) (T, error) {
	conf := getConfig()

	switch valueName {
	case "Count":
		if conf.Count < 1 {
			log.Print("Config.Count < 1. Defaulting to 1")
			return any(1).(T), nil
		}
		return any(conf.Count).(T), nil
	case "ProtocolICMP":
		if conf.ProtocolIMCP < 1 {
			log.Print("Config.ProtocolIMCP < 1. Defaulting to 1")
			return any(1).(T), nil
		}
		return any(conf.ProtocolIMCP).(T), nil
	case "TimeoutSecond":
		if conf.TimeoutSecond < 1 {
			log.Print("Config.TimeoutSecond < 1. Defaulting to 1")
			return any(1).(T), nil
		}
		return any(conf.TimeoutSecond).(T), nil
	default:
		var zero T
		return zero, fmt.Errorf("unknown config value: %s", valueName)
	}
}

func getConfig() *config.Config {
	config, err := config.NewConfig("config/config.json")
	if err != nil {
		log.Fatalf("failed to read config file. caused by : %v", err)
	}
	return config
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
