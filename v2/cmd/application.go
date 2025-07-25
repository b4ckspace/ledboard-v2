package cmd

import (
	"log/slog"
	"strconv"
	"time"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
	"github.com/b4ckspace/ledboard-v2/config"
	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/screens"
	"github.com/b4ckspace/ledboard-v2/utils"
)

// Mode represents the application's operational mode.
type Mode string

const (
	DefaultMode     Mode = "default"
	LasercutterMode Mode = "lasercutter"
)

// MQTTClient defines the interface for MQTT client operations.
type MQTTClient interface {
	Subscribe(topic string, handler mqttlib.MessageHandler) error
}

// RunApplication runs the application based on the specified mode.
func RunApplication(cfg *config.Config, ledBoardClient ledboard.LEDBoardClient, mqttClient MQTTClient, mode Mode) {
	// Initialize screens manager
	screensManager := screens.NewScreens()

	err := ledBoardClient.Init()
	if err != nil {
		slog.Error("Failed to initialize LED board client", "error", err)
		return
	}

	var memberCount int
	laserActive := false

	// Define idleScreen function based on laserActive status and mode
	getIdleScreen := func() string {
		if mode == LasercutterMode && laserActive {
			return screensManager.LaserOperation()
		}
		return screensManager.Idle(memberCount)
	}

	// Set time initially
	ledBoardClient.SetDate(time.Now())

	// Define a message handler for MQTT messages
	messageHandler := func(client mqttlib.Client, msg mqttlib.Message) {
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

			// Only send idle screen if not in laser mode or laser is not active
			if !laserActive || mode == DefaultMode {
				ledBoardClient.SendScreen(screensManager.Idle(memberCount))
			}

		case "psa/pizza":
			ledBoardClient.SendScreens([]string{screensManager.PizzaTimer(), getIdleScreen()})

		case "psa/donation":
			ledBoardClient.SendScreens([]string{screensManager.Donation(), getIdleScreen()})

		case "psa/alarm":
			ledBoardClient.SendScreens([]string{screensManager.Alarm(message), getIdleScreen()})

		case "psa/newMember":
			ledBoardClient.SendScreens([]string{screensManager.NewMemberRegistration(message), getIdleScreen()})

		case "sensor/door/bell":
			if message == "pressed" {
				ledBoardClient.SendScreens([]string{screensManager.DoorBell(), getIdleScreen()})
			}

		case "psa/message":
			if message != "" {
				ledBoardClient.SendScreens([]string{screensManager.PublicServiceAnnouncement(message), getIdleScreen()})
			}

		case "psa/nowPlaying":
			if message != "" {
				ledBoardClient.SendScreens([]string{screensManager.NowPlaying(message), getIdleScreen()})
			}

		// Laser-specific topics - these cases will only be reached if subscribed
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
		}
	}

	// Common MQTT subscriptions
	if err := mqttClient.Subscribe("psa/alarm", messageHandler); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/alarm", "error", err)
	}
	if err := mqttClient.Subscribe("psa/pizza", messageHandler); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/pizza", "error", err)
	}
	if err := mqttClient.Subscribe("psa/message", messageHandler); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/message", "error", err)
	}
	if err := mqttClient.Subscribe("sensor/door/bell", messageHandler); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "sensor/door/bell", "error", err)
	}
	if err := mqttClient.Subscribe("sensor/space/member/present", messageHandler); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "sensor/space/member/present", "error", err)
	}

	// Mode-specific MQTT subscriptions
	switch mode {
	case DefaultMode:
		if err := mqttClient.Subscribe("psa/donation", messageHandler); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/donation", "error", err)
		}
		if err := mqttClient.Subscribe("psa/newMember", messageHandler); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/newMember", "error", err)
		}
		if err := mqttClient.Subscribe("psa/nowPlaying", messageHandler); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/nowPlaying", "error", err)
		}
	case LasercutterMode:
		if err := mqttClient.Subscribe("project/laser/operation", messageHandler); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "project/laser/operation", "error", err)
		}
		if err := mqttClient.Subscribe("project/laser/finished", messageHandler); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "project/laser/finished", "error", err)
		}
		if err := mqttClient.Subscribe("project/laser/duration", messageHandler); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "project/laser/duration", "error", err)
		}
	}

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