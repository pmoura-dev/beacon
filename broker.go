package beacon

type Broker interface {
	Connect() error
	Disconnect() error

	Subscriber
	Publisher
}

type Subscriber interface {
	Subscribe(topic string) (<-chan Message, error)
}

type Publisher interface {
	Publish(topic string, message Message) error
}
