package application

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/b4ckspace/ledboard-v2/ledboard"
	"github.com/b4ckspace/ledboard-v2/mqttclient"
	"github.com/b4ckspace/ledboard-v2/screens"
	"github.com/b4ckspace/ledboard-v2/utils"

	mqttlib "github.com/eclipse/paho.mqtt.golang"
)

// Mode represents the application's operational mode.
type Mode string

const (
	DefaultMode     Mode = "default"
	LasercutterMode Mode = "lasercutter"
)

// Application holds the dependencies and state for the MQTT message handler.
type Application struct {
	ledBoardClient *ledboard.Client
	mqttClient     *mqttclient.Client
	pingProbe      *utils.PingProbe
	screens        *screens.Screens

	mode Mode

	memberCount int
	laserActive bool
}

// NewApplication creates a new Application.
func NewApplication(ledBoardClient *ledboard.Client, mqttClient *mqttclient.Client, pingProbe *utils.PingProbe, mode Mode) *Application {
	return &Application{
		ledBoardClient: ledBoardClient,
		mqttClient:     mqttClient,
		pingProbe:      pingProbe,
		screens:        screens.NewScreens(),
		mode:           mode,
	}
}

// getIdleScreen returns the appropriate idle screen based on current state.
func (app *Application) getIdleScreen() string {
	if app.mode == LasercutterMode && app.laserActive {
		return app.screens.LaserOperation()
	}
	return app.screens.Idle(app.memberCount)
}

// Run runs the application based on the specified mode.
func (app *Application) Run(ctx context.Context) error {
	// Common MQTT subscriptions
	if err := app.mqttClient.Subscribe("psa/alarm", app.handleMQTTMessage); err != nil {
		return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "psa/alarm", err)
	}
	if err := app.mqttClient.Subscribe("psa/pizza", app.handleMQTTMessage); err != nil {
		return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "psa/pizza", err)
	}
	if err := app.mqttClient.Subscribe("psa/message", app.handleMQTTMessage); err != nil {
		return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "psa/message", err)
	}
	if err := app.mqttClient.Subscribe("sensor/door/bell", app.handleMQTTMessage); err != nil {
		return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "sensor/door/bell", err)
	}
	if err := app.mqttClient.Subscribe("sensor/space/member/present", app.handleMQTTMessage); err != nil {
		return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "sensor/space/member/present", err)
	}

	// Mode-specific MQTT subscriptions
	switch app.mode {
	case DefaultMode:
		if err := app.mqttClient.Subscribe("psa/donation", app.handleMQTTMessage); err != nil {
			return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "psa/donation", err)
		}
		if err := app.mqttClient.Subscribe("psa/newMember", app.handleMQTTMessage); err != nil {
			return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "psa/newMember", err)
		}
		if err := app.mqttClient.Subscribe("psa/nowPlaying", app.handleMQTTMessage); err != nil {
			return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "psa/nowPlaying", err)
		}
	case LasercutterMode:
		if err := app.mqttClient.Subscribe("project/laser/operation", app.handleMQTTMessage); err != nil {
			return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "project/laser/operation", err)
		}
		if err := app.mqttClient.Subscribe("project/laser/finished", app.handleMQTTMessage); err != nil {
			return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "project/laser/finished", err)
		}
		if err := app.mqttClient.Subscribe("project/laser/duration", app.handleMQTTMessage); err != nil {
			return fmt.Errorf("failed to subscribe to MQTT topic: %s, %s", "project/laser/duration", err)
		}
	}

	err := app.pingProbe.Run(ctx, func() {
		slog.Info("ledboard is alive, setting date and sending idle screen")
		app.ledBoardClient.SetDate(time.Now())
		app.ledBoardClient.SendScreen(app.getIdleScreen())
	})
	if err != nil {
		return fmt.Errorf("issues while pinging: %s", err)
	}

	<-ctx.Done()
	slog.Info("Application context cancelled. Disconnecting MQTT client.")
	app.mqttClient.Disconnect()
	return nil
}

// handleMQTTMessage processes incoming MQTT messages.
func (app *Application) handleMQTTMessage(client mqttlib.Client, msg mqttlib.Message) {
	message := string(msg.Payload())
	slog.Info("Received MQTT message", "topic", msg.Topic(), "value", message)

	switch msg.Topic() {
	case "sensor/space/member/present":
		count, err := strconv.Atoi(message)
		if err != nil {
			slog.Error("Error converting member count", "error", err)
			return
		}
		app.memberCount = count

		if !app.laserActive {
			app.ledBoardClient.SendScreen(app.getIdleScreen())
		}

	case "psa/pizza":
		app.ledBoardClient.SendScreens([]string{app.screens.PizzaTimer(), app.getIdleScreen()})

	case "psa/donation":
		app.ledBoardClient.SendScreens([]string{app.screens.Donation(), app.getIdleScreen()})

	case "psa/alarm":
		app.ledBoardClient.SendScreens([]string{app.screens.Alarm(message), app.getIdleScreen()})

	case "psa/newMember":
		app.ledBoardClient.SendScreens([]string{app.screens.NewMemberRegistration(message), app.getIdleScreen()})

	case "sensor/door/bell":
		if message == "pressed" {
			app.ledBoardClient.SendScreens([]string{app.screens.DoorBell(), app.getIdleScreen()})
		}

	case "psa/message":
		if message != "" {
			app.ledBoardClient.SendScreens([]string{app.screens.PublicServiceAnnouncement(message), app.getIdleScreen()})
		}

	case "psa/nowPlaying":
		if message != "" {
			app.ledBoardClient.SendScreens([]string{app.screens.NowPlaying(message), app.getIdleScreen()})
		}

	case "project/laser/operation":
		if message == "active" {
			app.laserActive = true

			// Use the internal datetime to produce a counting screen!
			nullDate := time.Date(2000, time.February, 0, 0, 0, 2, 0, time.UTC)
			app.ledBoardClient.SetDate(nullDate)

			app.ledBoardClient.SendScreen(app.screens.LaserOperation())
		} else {
			app.laserActive = false
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
			app.ledBoardClient.SetDate(correction)
		}

	case "project/laser/finished":
		if message != "" {
			duration, err := strconv.Atoi(message)
			if err != nil {
				slog.Error("Error converting duration", "error", err)
				return
			}
			app.ledBoardClient.SendScreens([]string{app.screens.LaserFinished(duration), app.getIdleScreen()})

			// Reset datetime to something useful
			app.ledBoardClient.SetDate(time.Now())
		}
	}
}
