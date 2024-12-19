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
	UploadURLs   []string `json"uploadUrls"`
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
			"http://ipv4.download.thinkbroadband.com/512MB.zip",
			"http://ipv4.download.thinkbroadband.com/200MB.zip",
			"http://ipv4.download.thinkbroadband.com/100MB.zip",
		}
	}
	if len(config.Tests.SpeedTestURLs.UploadURLs) == 0 {
		config.Tests.SpeedTestURLs.DownloadURLs = []string{
			"https://httpbin.org/post",
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
