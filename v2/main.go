package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/b4ckspace/ledboard-v2/cmd"
	"github.com/b4ckspace/ledboard-v2/config"
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/mqttclient"
)

func main() {
	configPath := flag.String("config", "", "Full path to the configuration JSON file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Set up a default logger
	var logLevel slog.Level
	if *debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))

	if *configPath == "" {
		slog.Error("Usage: go run main.go --config <full_path_to_config.json>")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	slog.Info("Loaded configuration", "config", fmt.Sprintf("%+v", cfg))

	// Initialize LED Board Client (concrete implementation)
	ledBoardClient := ledboard.NewClient(cfg.LedBoardHost, 9520)

	// Connect to MQTT
	err = mqttclient.Connect(cfg)
	if err != nil {
		slog.Error("Failed to connect to MQTT broker", "error", err)
		os.Exit(1)
	}
	defer mqttclient.Disconnect()

	switch cfg.Mode {
	case "default":
		cmd.RunDefaultMode(cfg, ledBoardClient)
	case "lasercutter":
		cmd.RunLasercutterMode(cfg, ledBoardClient)
	default:
		slog.Error("Unknown configuration mode", "mode", cfg.Mode)
		os.Exit(1)
	}
}