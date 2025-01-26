package beacon

type Broker interface {
	Connect() error
	Disconnect() error

	Subscribe(topic string) (<-chan Message, error)
	Publish(topic string, message Message) error
}
