package cmd

import (
	"log/slog"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/b4ckspace/ledboard-v2/config"
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/mqttclient"
	"github.com/b4ckspace/ledboard-v2/screens"
	"github.com/b4ckspace/ledboard-v2/utils"
)

// RunLasercutterMode runs the application in lasercutter mode.
func RunLasercutterMode(cfg *config.Config) {
	// Initialize screens manager
	screensManager := screens.NewScreens()

	// Initialize LED Board Client
	ledBoardClient := ledboard.NewClient(cfg.LedBoardHost, 9520)
	err := ledBoardClient.Init()
	if err != nil {
		slog.Error("Failed to initialize LED board client", "error", err)
		return
	}

	var memberCount int
	laserActive := false

	// Define idleScreen function based on laserActive status
	getIdleScreen := func() string {
		if laserActive {
			return screensManager.LaserOperation()
		}
		return screensManager.Idle(memberCount)
	}

	// Set time initially
	ledBoardClient.SetDate(time.Now())

	// Define a message handler for MQTT messages
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		message := string(msg.Payload())
		slog.Info("Received MQTT message", "topic", msg.Topic(), "value", message)

		switch msg.Topic() {
		case "sensor/space/member/present":
			count, err := strconv.Atoi(message)
			if err != nil {
				slog.Error("Error converting member count", "error", err)
				return
			}
			memberCount = count

			if !laserActive {
				ledBoardClient.SendScreen(screensManager.Idle(memberCount))
			}

		case "project/laser/operation":
			if message == "active" {
				laserActive = true

				// Use the internal datetime to produce a counting screen!
				nullDate := time.Date(2000, time.February, 0, 0, 0, 2, 0, time.UTC)
				ledBoardClient.SetDate(nullDate)

				ledBoardClient.SendScreen(screensManager.LaserOperation())
			} else {
				laserActive = false
			}

		case "project/laser/duration":
			duration, err := strconv.Atoi(message)
			if err != nil {
				slog.Error("Error converting duration", "error", err)
				return
			}

			minutes := (duration % 3600) / 60
			seconds := duration % 60

			if minutes%2 == 0 && seconds == 57 {
				correction := time.Date(2000, time.February, 0, 0, minutes+1, 0, 0, time.UTC)
				ledBoardClient.SetDate(correction)
			}

		case "project/laser/finished":
			if message != "" {
				duration, err := strconv.Atoi(message)
				if err != nil {
					slog.Error("Error converting duration", "error", err)
					return
				}
				ledBoardClient.SendScreens([]string{screensManager.LaserFinished(duration), getIdleScreen()})

				// Reset datetime to something useful
				ledBoardClient.SetDate(time.Now())
			}

		case "psa/pizza":
			ledBoardClient.SendScreens([]string{screensManager.PizzaTimer(), getIdleScreen()})

		case "psa/alarm":
			ledBoardClient.SendScreens([]string{screensManager.Alarm(message), getIdleScreen()})

		case "sensor/door/bell":
			if message == "pressed" {
				ledBoardClient.SendScreens([]string{screensManager.DoorBell(), getIdleScreen()})
			}

		case "psa/message":
			if message != "" {
				ledBoardClient.SendScreens([]string{screensManager.PublicServiceAnnouncement(message), getIdleScreen()})
			}
		}
	}

	// Subscribe to MQTT topics
	mqttclient.Subscribe("project/laser/operation", messageHandler)
	mqttclient.Subscribe("project/laser/finished", messageHandler)
	mqttclient.Subscribe("project/laser/duration", messageHandler)

	mqttclient.Subscribe("psa/alarm", messageHandler)
	mqttclient.Subscribe("psa/pizza", messageHandler)
	mqttclient.Subscribe("psa/message", messageHandler)
	mqttclient.Subscribe("sensor/door/bell", messageHandler)
	mqttclient.Subscribe("sensor/space/member/present", messageHandler)

	// PingProbe
	aliveProbe := utils.NewPingProbe(cfg.Mqtt.Host, cfg.Ping) // Assuming config.mqtt.host is the host to ping
	go aliveProbe.Start()

	go func() {
		for range aliveProbe.AliveEvents() {
			slog.Info("Host is alive! Setting date and sending idle screen.")
			ledBoardClient.SetDate(time.Now())
			ledBoardClient.SendScreen(getIdleScreen())
		}
	}()

	// Keep the main goroutine alive
	select {}
}