package brokers

import (
	"sync"

	"github.com/pmoura-dev/beacon"
)

type LocalBroker struct {
	subscriptions map[string][]chan beacon.Message

	mu     sync.Mutex
	closed bool
}

func NewLocalBroker() *LocalBroker {
	return &LocalBroker{
		subscriptions: make(map[string][]chan beacon.Message),
	}
}

func (b *LocalBroker) Subscribe(topic string) (<-chan beacon.Message, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil, nil
	}

	messageChan := make(chan beacon.Message)
	b.subscriptions[topic] = append(b.subscriptions[topic], messageChan)
	return messageChan, nil
}

func (b *LocalBroker) Publish(topic string, message beacon.Message) {
	for _, ch := range b.subscriptions[topic] {
		ch <- message
	}
}

func (b *LocalBroker) Close() {}
