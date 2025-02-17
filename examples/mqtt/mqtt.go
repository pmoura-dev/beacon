package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/pmoura-dev/beacon"
	"github.com/pmoura-dev/beacon/brokers"
)

func fooHandler(publisher beacon.Publisher, message beacon.RoutedMessage) error {

	fmt.Printf("message: %s\n", string(message.Payload))
	fmt.Printf("topic: %s\n", message.Topic.FullName())
	fmt.Printf("param foo_id: %s\n", message.GetTopicParam("foo_id"))

	pubTopic, _ := beacon.NewTopic("bar/topic")
	err := publisher.Publish(pubTopic, beacon.Message{Payload: message.Payload})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	mqttURL := "tcp://broker.hivemq.com:1883"

	r := beacon.NewRouter(
		brokers.NewMQTTBroker(mqttURL),
	)

	_ = r.AddSubscription("foo/{foo_id}/topic", fooHandler)

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
