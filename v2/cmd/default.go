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

// RunDefaultMode runs the application in default mode.
func RunDefaultMode(cfg *config.Config, ledBoardClient ledboard.LEDBoardClient) {
	// Initialize screens manager
	screensManager := screens.NewScreens()

	err := ledBoardClient.Init()
	if err != nil {
		slog.Error("Failed to initialize LED board client", "error", err)
		return
	}

	var memberCount int

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
			ledBoardClient.SendScreen(screensManager.Idle(memberCount))

		case "psa/pizza":
			ledBoardClient.SendScreens([]string{screensManager.PizzaTimer(), screensManager.Idle(memberCount)})

		case "psa/donation":
			ledBoardClient.SendScreens([]string{screensManager.Donation(), screensManager.Idle(memberCount)})

		case "psa/alarm":
			ledBoardClient.SendScreens([]string{screensManager.Alarm(message), screensManager.Idle(memberCount)})

		case "psa/newMember":
			ledBoardClient.SendScreens([]string{screensManager.NewMemberRegistration(message), screensManager.Idle(memberCount)})

		case "sensor/door/bell":
			if message == "pressed" {
				ledBoardClient.SendScreens([]string{screensManager.DoorBell(), screensManager.Idle(memberCount)})
			}

		case "psa/message":
			if message != "" {
				ledBoardClient.SendScreens([]string{screensManager.PublicServiceAnnouncement(message), screensManager.Idle(memberCount)})
			}

		case "psa/nowPlaying":
			if message != "" {
				ledBoardClient.SendScreens([]string{screensManager.NowPlaying(message), screensManager.Idle(memberCount)})
			}
		}
	}

	// Subscribe to MQTT topics
	mqttclient.Subscribe("psa/alarm", messageHandler)
	mqttclient.Subscribe("psa/donation", messageHandler)
	mqttclient.Subscribe("psa/pizza", messageHandler)
	mqttclient.Subscribe("psa/newMember", messageHandler)
	mqttclient.Subscribe("psa/message", messageHandler)
	mqttclient.Subscribe("psa/nowPlaying", messageHandler)
	mqttclient.Subscribe("sensor/door/bell", messageHandler)
	mqttclient.Subscribe("sensor/space/member/present", messageHandler)

	// PingProbe
	aliveProbe := utils.NewPingProbe(cfg.Mqtt.Host, cfg.Ping) // Assuming config.mqtt.host is the host to ping
	go aliveProbe.Start()

	go func() {
		for range aliveProbe.AliveEvents() {
			slog.Info("Host is alive! Setting date and sending idle screen.")
			ledBoardClient.SetDate(time.Now())
			ledBoardClient.SendScreen(screensManager.Idle(memberCount))
		}
	}()

	// Keep the main goroutine alive
	select {}
}
