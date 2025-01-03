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
	ICMP          ICMPConfig      `json:"icmp"`
	SpeedTestURLs SpeedTestURLs   `json:"speedTestURLs"`
	RouteTest     RouteConfig     `json:"routeTest"`
	JitterTest    JitterConfig    `json:"jitterTest"`
	Bandwidth     BandwidthConfig `json:"bandwidth"`
}

type ICMPConfig struct {
	PacketCount    int `json:"packetCount"`
	TimeoutSeconds int `json:"timeoutSeconds"`
}

type SpeedTestURLs struct {
	DownloadURLs []string `json"downloadUrls"`
	UploadURLs   []string `json"uploadUrls"`
}

type RouteConfig struct {
	Target         string `json:"target"`
	MaxHops        int    `json:"maxHops"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
}

type JitterConfig struct {
	Target         string `json:"target"`
	PacketCount    int    `json:"packetCount"`
	TimeoutSeconds int    `json:"timeoutSeconds"`
}

type BandwidthConfig struct {
	InitialConnections int     `json:"initialConnections"`
	MaxConnections     int     `json:"maxConnections"`
	StepSize           int     `json:"stepSize"`
	FailThreshold      float64 `json:"failThreshold"`
	DownloadURL        string  `json"downloadUrl"`
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
			"https://httpbin.org/anything",
			"https://catbox.moe",
		}
	}
	if config.Tests.RouteTest.MaxHops == 0 {
		config.Tests.RouteTest.MaxHops = 30
	}
	if config.Tests.RouteTest.TimeoutSeconds == 0 {
		config.Tests.RouteTest.TimeoutSeconds = 5
	}
	if config.Tests.JitterTest.PacketCount == 0 {
		config.Tests.JitterTest.PacketCount = 10
	}
	if config.Tests.JitterTest.TimeoutSeconds == 0 {
		config.Tests.JitterTest.TimeoutSeconds = 5
	}

	if config.Tests.Bandwidth.InitialConnections == 0 {
		config.Tests.Bandwidth.InitialConnections = 1
	}
	if config.Tests.Bandwidth.MaxConnections == 0 {
		config.Tests.Bandwidth.MaxConnections = 32
	}
	if config.Tests.Bandwidth.StepSize == 0 {
		config.Tests.Bandwidth.StepSize = 2
	}
	if config.Tests.Bandwidth.FailThreshold == 0 {
		config.Tests.Bandwidth.FailThreshold = 70 // bandwidth falls off fast with most isps
	}
	if len(config.Tests.Bandwidth.DownloadURL) == 0 {
		config.Tests.Bandwidth.DownloadURL = "http://ipv4.download.thinkbroadband.com/100MB.zip"
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
