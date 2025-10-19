package subscribers

import (
	"regexp"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pmoura-dev/beacon"
)

type MQTTSubscriber struct {
	client               mqtt.Client
	qos                  byte
	disconnectionTimeout uint // milliseconds
}

type MQTTSubscriberOption func(*MQTTSubscriber)

func NewMQTTSubscriber(url string, options ...MQTTSubscriberOption) *MQTTSubscriber {
	opts := mqtt.NewClientOptions().AddBroker(url)
	opts.SetCleanSession(true)

	subscriber := &MQTTSubscriber{
		client:               mqtt.NewClient(opts),
		qos:                  0,
		disconnectionTimeout: 250,
	}

	for _, opt := range options {
		opt(subscriber)
	}

	return subscriber
}

func WithQOS(qos byte) func(*MQTTSubscriber) {
	return func(b *MQTTSubscriber) {
		b.qos = qos
	}
}

func WithDisconnectionTimeout(timeout uint) func(*MQTTSubscriber) {
	return func(b *MQTTSubscriber) {
		b.disconnectionTimeout = timeout
	}
}

func (b *MQTTSubscriber) Connect() error {
	if token := b.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (b *MQTTSubscriber) Disconnect() error {
	b.client.Disconnect(b.disconnectionTimeout)
	return nil
}

func (b *MQTTSubscriber) Subscribe(topic *beacon.Topic) (<-chan beacon.RoutedMessage, error) {
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
