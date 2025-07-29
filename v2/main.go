package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/b4ckspace/ledboard-v2/cmd"
	"github.com/b4ckspace/ledboard-v2/config"
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/mqttclient"
	"github.com/b4ckspace/ledboard-v2/utils"
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

	// Initialize MQTT Client
	mqttClient := mqttclient.NewClient()
	err = mqttClient.Connect(cfg)
	if err != nil {
		slog.Error("Failed to connect to MQTT broker", "error", err)
		os.Exit(1)
	}
	// Defer mqttClient.Disconnect() is moved to application.go to ensure it's called on graceful shutdown

	// Initialize PingProbe
	aliveProbe := utils.NewPingProbe(cfg.LedBoardHost, cfg.Ping)

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called when main exits

	// Listen for OS signals to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		slog.Info("Received signal, shutting down...", "signal", sig)
		cancel() // Trigger context cancellation
	}()

	var app *cmd.Application
	switch cfg.Mode {
	case string(cmd.DefaultMode):
		app = cmd.NewApplication(cfg, ledBoardClient, mqttClient, aliveProbe, cmd.DefaultMode, ctx)
	case string(cmd.LasercutterMode):
		app = cmd.NewApplication(cfg, ledBoardClient, mqttClient, aliveProbe, cmd.LasercutterMode, ctx)
	default:
		slog.Error("Unknown configuration mode", "mode", cfg.Mode)
		os.Exit(1)
	}

	app.Run()
}
