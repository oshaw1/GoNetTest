package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port          int    `json:"port"`
	Host          string `json:"host"`
	Count         int    `json:"ICMP packets"`
	ProtocolIMCP  int    `json:"protocol IMCP"`
	TimeoutSecond int    `json:"Timeout Threshold (Seconds)"`
}

func NewConfig(filepath string) (*Config, error) {
	config, err := load(filepath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func load(filename string) (*Config, error) {
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
