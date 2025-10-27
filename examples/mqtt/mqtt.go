package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"os/signal"
	"syscall"
	"time"

	"github.com/pmoura-dev/beacon"
	"github.com/pmoura-dev/beacon/publishers"
	"github.com/pmoura-dev/beacon/subscribers"
)

func fooHandler(ctx context.Context, publisher beacon.Publisher, message beacon.RoutedMessage) error {

	fmt.Printf("message: %s\n", string(message.Payload))
	fmt.Printf("topic: %s\n", message.Topic.FullName())
	fmt.Printf("param foo_id: %s\n", message.GetTopicParam("foo_id"))
	fmt.Printf("context value: %d\n", ctx.Value("c_value").(int))

	pubTopic, _ := beacon.NewTopic("bar/topic")
	err := publisher.Publish(pubTopic, beacon.Message{Payload: message.Payload})
	if err != nil {
		return err
	}

	return nil
}

func timingMiddleware(next beacon.HandlerFunc) beacon.HandlerFunc {
	return func(ctx context.Context, publisher beacon.Publisher, message beacon.RoutedMessage) error {
		startTime := time.Now()

		next(ctx, publisher, message)

		elapsedTime := time.Since(startTime)
		log.Printf("[%s] [%s]\n", message.Topic.FullName(), elapsedTime)
		return nil
	}
}

func injectMiddleware(next beacon.HandlerFunc) beacon.HandlerFunc {
	return func(ctx context.Context, publisher beacon.Publisher, message beacon.RoutedMessage) error {
		x := rand.IntN(100)
		ctx = context.WithValue(ctx, "c_value", x)

		next(ctx, publisher, message)

		return nil
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	mqttURL := "mqtt://broker.emqx.io:1883"

	subscriber := subscribers.NewMQTTSubscriber(mqttURL)
	publisher := publishers.NewMQTTPublisher(mqttURL)

	r := beacon.NewRouter(
		beacon.NewBroker(subscriber, publisher),
	)

	_ = r.AddSubscription("foo/{foo_id}/topic", fooHandler)

	_ = r.UseMiddleware(timingMiddleware)
	_ = r.UseMiddleware(injectMiddleware)

	if err := r.Start(); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := r.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Error shutting down Beacon.")
	}
}
