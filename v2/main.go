package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/b4ckspace/ledboard-v2/application"
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/mqttclient"
	"github.com/b4ckspace/ledboard-v2/utils"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug bool `envconfig:"DEBUG"`

	Mode         string `envconfig:"MODE" required:"true"`
	LedBoardHost string `envconfig:"LEDBOARD_HOST" required:"true"`

	LedBoardPingIntervalSeconds int `envconfig:"LEDBOARD_PING_INTERVAL_SECONDS" required:"true"`

	MqttHost string `envconfig:"MQTT_HOST" required:"true"`
}

func main() {

	config := Config{}
	err := envconfig.Process("", &config)
	if err != nil {
		slog.Error("unable to parse environment", "error", err)
		os.Exit(1)
	}

	// Set up a default logger
	logLevel := slog.LevelInfo
	if config.Debug {
		logLevel = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))

	// Initialize LED Board Client
	ledBoardClient, err := ledboard.NewClient(config.LedBoardHost, 9520)
	if err != nil {
		slog.Error("failed to connect to ledboard", "error", err)
		os.Exit(1)
	}

	// Initialize MQTT Client
	mqttClient := mqttclient.NewClient()
	err = mqttClient.Connect(config.MqttHost)
	if err != nil {
		slog.Error("failed to connect to mqtt broker", "error", err)
		os.Exit(1)
	}
	defer mqttClient.Disconnect()

	// Set up context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize PingProbe
	pingProbe, err := utils.NewPingProbe(config.LedBoardHost, config.LedBoardPingIntervalSeconds)
	if err != nil {
		slog.Error("unable to start pingprobe", "error", err)
		os.Exit(1)
	}

	// Listen for OS signals to gracefully shut down
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		slog.Info("received signal, shutting down...", "signal", sig)
		cancel()
	}()

	var app *application.Application
	switch config.Mode {
	case string(application.DefaultMode):
		fallthrough
	case string(application.LasercutterMode):
		app = application.NewApplication(ledBoardClient, mqttClient, pingProbe, application.Mode(config.Mode))
	default:
		slog.Error("unknown configuration mode", "mode", config.Mode)
		os.Exit(1)
	}

	err = app.Run(ctx)
	if err != nil {
		slog.Error("Unable to run", "error", err)
		os.Exit(1)
	}
}
