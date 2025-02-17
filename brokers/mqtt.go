package brokers

import (
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pmoura-dev/beacon"
)

type MQTTBroker struct {
	client               mqtt.Client
	qos                  byte
	disconnectionTimeout uint // milliseconds
}

type MQTTBrokerOption func(*MQTTBroker)

func NewMQTTBroker(url string, options ...MQTTBrokerOption) *MQTTBroker {
	opts := mqtt.NewClientOptions().AddBroker(url)
	opts.SetCleanSession(true)

	broker := &MQTTBroker{
		client:               mqtt.NewClient(opts),
		qos:                  0,
		disconnectionTimeout: 250,
	}

	for _, opt := range options {
		opt(broker)
	}

	return broker
}

func WithQOS(qos byte) func(*MQTTBroker) {
	return func(b *MQTTBroker) {
		b.qos = qos
	}
}

func WithDisconnectionTimeout(timeout uint) func(*MQTTBroker) {
	return func(b *MQTTBroker) {
		b.disconnectionTimeout = timeout
	}
}

func (b *MQTTBroker) Connect() error {
	if token := b.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (b *MQTTBroker) Disconnect() error {
	b.client.Disconnect(b.disconnectionTimeout)
	return nil
}

func (b *MQTTBroker) Subscribe(topic *beacon.Topic) (<-chan beacon.RoutedMessage, error) {
	messageChan := make(chan beacon.RoutedMessage)

	mqttTopic := toMQTTTopic(topic.Raw())
	token := b.client.Subscribe(mqttTopic, b.qos, func(c mqtt.Client, m mqtt.Message) {
		topicMatch := extractParamsFromMQTTTopic(topic, m.Topic())
		messageChan <- beacon.RoutedMessage{
			Message: beacon.Message{
				Payload: m.Payload(),
			},
			Topic: topicMatch,
		}
	})
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return messageChan, nil
}

func (b *MQTTBroker) Publish(topic *beacon.Topic, message beacon.Message) error {
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

func extractParamsFromMQTTTopic(topic *beacon.Topic, mqttTopic string) *beacon.TopicMatch {
	// {param} -> ([^/]+)
	rePattern := regexp.MustCompile(`\{([^}]+)\}`).ReplaceAllString(topic.Raw(), `([^/]+)`)
	rePattern = strings.ReplaceAll(rePattern, "*", ".*")

	re := regexp.MustCompile("^" + rePattern + "$")

	groups := re.FindStringSubmatch(mqttTopic)
	params := map[string]string{}
	for i, g := range groups[1:] {
		params[topic.Params()[i]] = g
	}

	return beacon.NewTopicMatch(mqttTopic, params)
}
