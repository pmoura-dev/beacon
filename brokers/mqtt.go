package brokers

import (
	"github.com/pmoura-dev/beacon"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTBroker struct {
	client mqtt.Client
}

func NewMQTTBroker(url string) *MQTTBroker {
	opts := mqtt.NewClientOptions().AddBroker(url)
	opts.SetCleanSession(true)

	return &MQTTBroker{
		client: mqtt.NewClient(opts),
	}
}

func (b *MQTTBroker) Connect() error {
	if token := b.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (b *MQTTBroker) Disconnect() error {
	b.client.Disconnect(250)
	return nil
}

func (b *MQTTBroker) Subscribe(topic string) (<-chan beacon.Message, error) {
	messageChan := make(chan beacon.Message)

	token := b.client.Subscribe(topic, 0, func(c mqtt.Client, m mqtt.Message) {
		message := b.toBeaconMessage(m)
		messageChan <- message
	})
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return messageChan, nil
}

func (b *MQTTBroker) Publish(topic string, message beacon.Message) error {
	return nil
}

func (b *MQTTBroker) Close() {
	b.client.Disconnect(250)
}

func (b *MQTTBroker) toBeaconMessage(mqttMessage mqtt.Message) beacon.Message {
	return beacon.Message{
		Payload: mqttMessage.Payload(),
	}
}
