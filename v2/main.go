package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/b4ckspace/ledboard-v2/cmd"
	"github.com/b4ckspace/ledboard-v2/config"
	"github.com/b4ckspace/ledboard-v2/mqttclient"
)

func main() {
	configPath := flag.String("config", "", "Full path to the configuration JSON file")
	flag.Parse()

	if *configPath == "" {
		fmt.Println("Usage: go run main.go --config <full_path_to_config.json>")
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	fmt.Printf("Loaded configuration: %+v\n", cfg)

	// Connect to MQTT
	err = mqttclient.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttclient.Disconnect()

	switch cfg.Mode {
	case "default":
		cmd.RunDefaultMode(cfg)
	case "lasercutter":
		cmd.RunLasercutterMode(cfg)
	default:
		log.Fatalf("Unknown configuration mode: %s", cfg.Mode)
	}
}
