package beacon

type Broker interface {
	Connect() error
	Disconnect() error

	Subscriber
	Publisher
}

type Subscriber interface {
	Subscribe(topic *Topic) (<-chan RoutedMessage, error)
}

type Publisher interface {
	Publish(topic *Topic, message Message) error
}
