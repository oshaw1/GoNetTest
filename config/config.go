package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	// UI Settings
	Dash DashboardSettings `json:"dashboard"` // How many days of tests to show in UI

	// Test Settings
	Tests TestConfigs `json:"tests"`
}

type DashboardSettings struct {
	RecentDays int `json:"recentDays"`
}

type TestConfigs struct {
	ICMP          ICMPConfig    `json:"icmp"`
	SpeedTestURLs SpeedTestURLs `json:"speedTestURLs"`
}

type ICMPConfig struct {
	PacketCount    int `json:"packetCount"`
	TimeoutSeconds int `json:"timeoutSeconds"`
}

type SpeedTestURLs struct {
	DownloadURLs []string `json"downloadUrls"`
	//UploadURLs   []string `json"UploadUrls"`
}

func NewConfig(filepath string) (*Config, error) {
	config, err := load(filepath)
	if err != nil {
		return nil, err
	}

	if config.Dash.RecentDays == 0 {
		config.Dash.RecentDays = 7 // Default to showing last 7 days
	}

	if config.Tests.ICMP.PacketCount == 0 {
		config.Tests.ICMP.PacketCount = 4
	}
	if config.Tests.ICMP.TimeoutSeconds == 0 {
		config.Tests.ICMP.TimeoutSeconds = 5
	}
	// Set default speed test URLs if none provided
	if len(config.Tests.SpeedTestURLs.DownloadURLs) == 0 {
		config.Tests.SpeedTestURLs.DownloadURLs = []string{
			"https://speed.cloudflare.com/100MB",
			"https://storage.googleapis.com/speed-test-files/100MB.bin",
			"https://speedtest-sfo2.digitalocean.com/100mb.test",
		}
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
