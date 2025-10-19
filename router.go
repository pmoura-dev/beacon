package beacon

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

var (
	ErrCannotAddSubscription   = errors.New("subscription can not be added")
	ErrDuplicateSubscription   = errors.New("subscription already exists")
	ErrShutdownTimeoutExceeded = errors.New("shutdown timeout exceeded")
)

type Router struct {
	broker *Broker
	logger *slog.Logger

	subscriptions map[*Topic]HandlerFunc

	wg sync.WaitGroup

	// Channel that indicates if the router is in shutdownChan mode.
	shutdownChan chan struct{}

	// Flag that indicates if the router was already started or not.
	isRunning bool
}

func NewRouter(broker *Broker, options ...OptionFunc) *Router {
	r := &Router{
		broker:        broker,
		logger:        slog.Default(),
		subscriptions: make(map[*Topic]HandlerFunc),

		shutdownChan: make(chan struct{}),
	}

	for _, opt := range options {
		opt(r)
	}

	return r
}

type OptionFunc func(*Router)

func WithLogger(logger *slog.Logger) func(*Router) {
	return func(r *Router) {
		r.logger = logger
	}
}

func (r *Router) Start() error {
	r.logger.Info("Starting Beacon...")

	err := r.broker.Connect()
	if err != nil {
		return err
	}

	r.logger.Info("Connected to broker.")

	r.startListening()

	r.logger.Info("Beacon started.")
	return nil
}

func (r *Router) startListening() {
	for topic, handler := range r.subscriptions {
		messageChan, err := r.broker.Subscribe(topic)
		if err != nil {
			r.logger.Error("Error adding subscription", "topic", topic, "error", err)
			continue
		}

		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			for {
				select {
				case <-r.shutdownChan:
					r.logger.Info("Router is in shutdown phase. Stopped listening for messages.", "topic", topic)
					return
				case message := <-messageChan:
					err := handler(r.broker, message)
					if err != nil {
						r.logger.Error("Error processing message.", "error", err)
					}
				}
			}
		}()

		r.logger.Info("Added subscription.", "topic", topic)
	}
}

func (r *Router) Shutdown(ctx context.Context) error {
	r.logger.Info("Shutting down Beacon...")
	close(r.shutdownChan)

	_ = r.broker.Disconnect()

	r.logger.Info("Waiting for in-flight messages to be processed.")

	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		r.logger.Info("All messages were processed.")
	case <-ctx.Done():
		r.logger.Info("Forcing shutdown.", "error", ErrShutdownTimeoutExceeded)
	}

	r.logger.Info("Beacon shutdown.")
	return nil
}

func (r *Router) AddSubscription(rawTopic string, handler HandlerFunc) error {
	if r.isRunning {
		r.logger.Error("Subscription could not be added. Router is already running.", "topic", rawTopic)
		return ErrCannotAddSubscription
	}

	topic, err := NewTopic(rawTopic)
	if err != nil {
		r.logger.Error("Invalid topic definition.", "topic", rawTopic)
		return err
	}

	if _, exists := r.subscriptions[topic]; exists {
		r.logger.Error("A subscription to this topic already exists", "topic", rawTopic)
		return ErrDuplicateSubscription
	}

	r.subscriptions[topic] = handler
	return nil
}

func (r *Router) Publish(rawTopic string, message Message) error {
	topic, err := NewTopic(rawTopic)
	if err != nil {
		return err
	}

	return r.broker.Publish(topic, message)
}

type HandlerFunc func(Publisher, RoutedMessage) error
