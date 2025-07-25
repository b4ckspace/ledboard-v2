package mqttclient

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/b4ckspace/ledboard-v2/config"
)

var client mqtt.Client

func Connect(cfg *config.Config) error {
	opts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:1883", cfg.Mqtt.Host)).SetClientID("github.com/b4ckspace/ledboard-v2")

	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("MQTT Connection lost: %v", err)
	})
	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("MQTT Connected")
	})

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	return nil
}

func Publish(topic string, payload string) {
	if client == nil || !client.IsConnected() {
		log.Println("MQTT client not connected, cannot publish.")
		return
	}
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish message to topic %s: %v", topic, token.Error())
	}
}

func Subscribe(topic string, handler mqtt.MessageHandler) error {
	if client == nil || !client.IsConnected() {
		return fmt.Errorf("MQTT client not connected, cannot subscribe.")
	}
	token := client.Subscribe(topic, 0, handler)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	return nil
}

func Disconnect() {
	if client != nil && client.IsConnected() {
		client.Disconnect(250)
		log.Println("MQTT Disconnected")
	}
}
