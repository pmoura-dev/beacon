package beacon

import "errors"

var (
	ErrNoSubscriber = errors.New("broker does not have a subscriber associated")
	ErrNoPublisher  = errors.New("broker does not have a publisher associated")
)

type Broker struct {
	subscriber Subscriber
	publisher  Publisher
}

func NewBroker(subscriber Subscriber, publisher Publisher) *Broker {
	return &Broker{
		subscriber: subscriber,
		publisher:  publisher,
	}
}

func (b *Broker) Connect() error {
	if b.subscriber != nil {
		if err := b.subscriber.Connect(); err != nil {
			return err
		}
	}

	if b.publisher != nil {
		if err := b.publisher.Connect(); err != nil {
			return err
		}
	}

	return nil
}

func (b *Broker) Disconnect() error {
	if b.subscriber != nil {
		if err := b.subscriber.Disconnect(); err != nil {
			return err
		}
	}

	if b.publisher != nil {
		if err := b.publisher.Disconnect(); err != nil {
			return err
		}
	}

	return nil
}

func (b *Broker) Subscribe(topic *Topic) (<-chan RoutedMessage, error) {
	if b.subscriber == nil {
		return nil, ErrNoSubscriber
	}

	return b.subscriber.Subscribe(topic)
}

func (b *Broker) Publish(topic *Topic, message Message) error {
	if b.publisher == nil {
		return ErrNoPublisher
	}

	return b.publisher.Publish(topic, message)
}

type Connector interface {
	Connect() error
	Disconnect() error
}

type Subscriber interface {
	Connector
	Subscribe(topic *Topic) (<-chan RoutedMessage, error)
}

type Publisher interface {
	Connector
	Publish(topic *Topic, message Message) error
}
