package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	// UI Settings
	RecentDays int `json:"recentDays"` // How many days of tests to show in UI

	// Test Settings
	Tests TestConfigs `json:"tests"`
}

type TestConfigs struct {
	ICMP ICMPConfig `json:"icmp"`
	TCP  TCPConfig  `json:"tcp"`
}

type ICMPConfig struct {
	PacketCount    int `json:"packetCount"`
	TimeoutSeconds int `json:"timeoutSeconds"`
}

type TCPConfig struct {
	Ports          []int `json:"ports"`
	TimeoutSeconds int   `json:"timeoutSeconds"`
}

func NewConfig(filepath string) (*Config, error) {
	config, err := load(filepath)
	if err != nil {
		return nil, err
	}

	// Set defaults
	if config.RecentDays == 0 {
		config.RecentDays = 7 // Default to showing last 7 days
	}

	if config.Tests.ICMP.PacketCount == 0 {
		config.Tests.ICMP.PacketCount = 4
	}
	if config.Tests.ICMP.TimeoutSeconds == 0 {
		config.Tests.ICMP.TimeoutSeconds = 5
	}

	if len(config.Tests.TCP.Ports) == 0 {
		config.Tests.TCP.Ports = []int{80, 443} // Default to common ports
	}
	if config.Tests.TCP.TimeoutSeconds == 0 {
		config.Tests.TCP.TimeoutSeconds = 5
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
