package mqttclient

import (
	"fmt"
	"log/slog"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Client holds the MQTT client instance.
type Client struct {
	mqttClient mqtt.Client
}

// NewClient creates and returns a new MQTT Client instance.
func NewClient() *Client {
	return &Client{}
}

// Connect connects the MQTT client to the broker.
func (c *Client) Connect(host string) error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:1883", host))
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		slog.Warn("mqtt connection lost", "error", err)
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		slog.Info("mqtt connected")
	})

	c.mqttClient = mqtt.NewClient(opts)
	if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to mqtt broker: %w", token.Error())
	}

	return nil
}

// Publish publishes a message to the specified MQTT topic.
func (c *Client) Publish(topic string, payload string) {
	if c.mqttClient == nil || !c.mqttClient.IsConnected() {
		slog.Warn("MQTT client not connected, cannot publish")
		return
	}
	token := c.mqttClient.Publish(topic, 0, false, payload)
	token.Wait()
	if token.Error() != nil {
		slog.Error("Failed to publish message", "topic", topic, "error", token.Error())
	}
}

// Subscribe subscribes to the specified MQTT topic.
func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	if c.mqttClient == nil || !c.mqttClient.IsConnected() {
		return fmt.Errorf("MQTT client not connected, cannot subscribe")
	}
	token := c.mqttClient.Subscribe(topic, 0, handler)
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	return nil
}

// Disconnect disconnects the MQTT client from the broker.
func (c *Client) Disconnect() {
	if c.mqttClient != nil && c.mqttClient.IsConnected() {
		c.mqttClient.Disconnect(250)
		slog.Info("MQTT Disconnected")
	}
}
