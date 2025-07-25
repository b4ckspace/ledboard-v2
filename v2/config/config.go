package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
)

type StatusConfig struct {
	URL      string `json:"url"`
	Interval int    `json:"interval"`
}

type PingConfig struct {
	Interval           int `json:"interval"` // in seconds
	ConsecutiveAnswers int `json:"consecutiveAnswers"`
}

type MqttConfig struct {
	Host string `json:"host"`
}

type Config struct {
	Mode        string       `json:"mode"`
	LedBoardHost string       `json:"ledBoardHost"` // New field for LED board hostname
	Status      StatusConfig `json:"status"`
	Ping        PingConfig   `json:"ping"`	
	Mqtt        MqttConfig   `json:"mqtt"`
}

var activeConfig *Config

func LoadConfig(configPath string) (*Config, error) {
	slog.Info("Attempting to load config from", "path", configPath)

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config JSON: %w", err)
	}

	activeConfig = &cfg
	return activeConfig, nil
}

func GetConfig() *Config {
	return activeConfig
}