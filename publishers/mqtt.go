package publishers

import (
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pmoura-dev/beacon"
)

type MQTTPublisher struct {
	client               mqtt.Client
	qos                  byte
	disconnectionTimeout uint // milliseconds
}

type MQTTPublisherOption func(*MQTTPublisher)

func NewMQTTPublisher(url string, options ...MQTTPublisherOption) *MQTTPublisher {
	opts := mqtt.NewClientOptions().AddBroker(url)
	opts.SetCleanSession(true)

	publisher := &MQTTPublisher{
		client:               mqtt.NewClient(opts),
		qos:                  0,
		disconnectionTimeout: 250,
	}

	for _, opt := range options {
		opt(publisher)
	}

	return publisher
}

func WithQOS(qos byte) func(*MQTTPublisher) {
	return func(b *MQTTPublisher) {
		b.qos = qos
	}
}

func WithDisconnectionTimeout(timeout uint) func(*MQTTPublisher) {
	return func(b *MQTTPublisher) {
		b.disconnectionTimeout = timeout
	}
}

func (b *MQTTPublisher) Connect() error {
	if token := b.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (b *MQTTPublisher) Disconnect() error {
	b.client.Disconnect(b.disconnectionTimeout)
	return nil
}

func (b *MQTTPublisher) Publish(topic *beacon.Topic, message beacon.Message) error {
	mqttTopic := toMQTTTopic(topic.Raw())

	token := b.client.Publish(mqttTopic, b.qos, false, message.Payload)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func toMQTTTopic(topic string) string {
	// Replace named wildcards {param} with "+"
	mqttTopic := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(topic, "+")

	// Replace * with "#"
	mqttTopic = strings.ReplaceAll(mqttTopic, "*", "#")

	return mqttTopic
}
