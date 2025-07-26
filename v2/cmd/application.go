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

// ApplicationContext holds the dependencies and state for the MQTT message handler.
type ApplicationContext struct {
	cfg            *config.Config
	ledBoardClient ledboard.LEDBoardClient
	mqttClient     MQTTClient
	aliveProbe     utils.PingProbe
	screensManager *screens.Screens
	memberCount    int
	laserActive    bool
	mode           Mode
}

// NewApplicationContext creates a new ApplicationContext.
func NewApplicationContext(cfg *config.Config, ledBoardClient ledboard.LEDBoardClient, mqttClient MQTTClient, aliveProbe utils.PingProbe, mode Mode) *ApplicationContext {
	return &ApplicationContext{
		cfg:            cfg,
		ledBoardClient: ledBoardClient,
		mqttClient:     mqttClient,
		aliveProbe:     aliveProbe,
		screensManager: screens.NewScreens(), // Initialize screensManager here
		mode:           mode,
	}
}

// getIdleScreen returns the appropriate idle screen based on current state.
func (appCtx *ApplicationContext) getIdleScreen() string {
	if appCtx.mode == LasercutterMode && appCtx.laserActive {
		return appCtx.screensManager.LaserOperation()
	}
	return appCtx.screensManager.Idle(appCtx.memberCount)
}

// handleMQTTMessage processes incoming MQTT messages.
func (appCtx *ApplicationContext) handleMQTTMessage(client mqttlib.Client, msg mqttlib.Message) {
	message := string(msg.Payload())
	slog.Info("Received MQTT message", "topic", msg.Topic(), "value", message)

	switch msg.Topic() {
	case "sensor/space/member/present":
		count, err := strconv.Atoi(message)
		if err != nil {
			slog.Error("Error converting member count", "error", err)
			return
		}
		appCtx.memberCount = count

		// Only send idle screen if not in laser mode or laser is not active
		if !appCtx.laserActive || appCtx.mode == DefaultMode {
			appCtx.ledBoardClient.SendScreen(appCtx.screensManager.Idle(appCtx.memberCount))
		}

	case "psa/pizza":
		appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.PizzaTimer(), appCtx.getIdleScreen()})

	case "psa/donation":
		appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.Donation(), appCtx.getIdleScreen()})

	case "psa/alarm":
		appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.Alarm(message), appCtx.getIdleScreen()})

	case "psa/newMember":
		appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.NewMemberRegistration(message), appCtx.getIdleScreen()})

	case "sensor/door/bell":
		if message == "pressed" {
			appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.DoorBell(), appCtx.getIdleScreen()})
		}

	case "psa/message":
		if message != "" {
			appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.PublicServiceAnnouncement(message), appCtx.getIdleScreen()})
		}

	case "psa/nowPlaying":
		if message != "" {
			appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.NowPlaying(message), appCtx.getIdleScreen()})
		}

	// Laser-specific topics - these cases will only be reached if subscribed
	case "project/laser/operation":
		if message == "active" {
			appCtx.laserActive = true

			// Use the internal datetime to produce a counting screen!
			nullDate := time.Date(2000, time.February, 0, 0, 0, 2, 0, time.UTC)
			appCtx.ledBoardClient.SetDate(nullDate)

			appCtx.ledBoardClient.SendScreen(appCtx.screensManager.LaserOperation())
		} else {
			appCtx.laserActive = false
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
			appCtx.ledBoardClient.SetDate(correction)
		}

	case "project/laser/finished":
		if message != "" {
			duration, err := strconv.Atoi(message)
			if err != nil {
				slog.Error("Error converting duration", "error", err)
				return
			}
			appCtx.ledBoardClient.SendScreens([]string{appCtx.screensManager.LaserFinished(duration), appCtx.getIdleScreen()})

			// Reset datetime to something useful
			appCtx.ledBoardClient.SetDate(time.Now())
		}
	}
}

// Run runs the application based on the specified mode.
func (appCtx *ApplicationContext) Run() {
	err := appCtx.ledBoardClient.Init()
	if err != nil {
		slog.Error("Failed to initialize LED board client", "error", err)
		return
	}

	// Set time initially
	appCtx.ledBoardClient.SetDate(time.Now())

	// Common MQTT subscriptions
	if err := appCtx.mqttClient.Subscribe("psa/alarm", appCtx.handleMQTTMessage); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/alarm", "error", err)
	}
	if err := appCtx.mqttClient.Subscribe("psa/pizza", appCtx.handleMQTTMessage); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/pizza", "error", err)
	}
	if err := appCtx.mqttClient.Subscribe("psa/message", appCtx.handleMQTTMessage); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/message", "error", err)
	}
	if err := appCtx.mqttClient.Subscribe("sensor/door/bell", appCtx.handleMQTTMessage); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "sensor/door/bell", "error", err)
	}
	if err := appCtx.mqttClient.Subscribe("sensor/space/member/present", appCtx.handleMQTTMessage); err != nil {
		slog.Error("Failed to subscribe to MQTT topic", "topic", "sensor/space/member/present", "error", err)
	}

	// Mode-specific MQTT subscriptions
	switch appCtx.mode {
	case DefaultMode:
		if err := appCtx.mqttClient.Subscribe("psa/donation", appCtx.handleMQTTMessage); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/donation", "error", err)
		}
		if err := appCtx.mqttClient.Subscribe("psa/newMember", appCtx.handleMQTTMessage); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/newMember", "error", err)
		}
		if err := appCtx.mqttClient.Subscribe("psa/nowPlaying", appCtx.handleMQTTMessage); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "psa/nowPlaying", "error", err)
		}
	case LasercutterMode:
		if err := appCtx.mqttClient.Subscribe("project/laser/operation", appCtx.handleMQTTMessage); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "project/laser/operation", "error", err)
		}
		if err := appCtx.mqttClient.Subscribe("project/laser/finished", appCtx.handleMQTTMessage); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "project/laser/finished", "error", err)
		}
		if err := appCtx.mqttClient.Subscribe("project/laser/duration", appCtx.handleMQTTMessage); err != nil {
			slog.Error("Failed to subscribe to MQTT topic", "topic", "project/laser/duration", "error", err)
		}
	}

	// PingProbe
	go appCtx.aliveProbe.Start()

	go func() {
		for range appCtx.aliveProbe.AliveEvents() {
			slog.Info("Host is alive! Setting date and sending idle screen.")
			appCtx.ledBoardClient.SetDate(time.Now())
			appCtx.ledBoardClient.SendScreen(appCtx.getIdleScreen())
		}
	}()

	// Keep the main goroutine alive
	select {}
}
