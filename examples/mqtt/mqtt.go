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

func fooHandler(publisher beacon.Publisher, message beacon.Message) error {
	fmt.Println("received in foo_topic message: ", string(message.Payload))

	fmt.Println("publishing message in bar_topic ")
	err := publisher.Publish("bar_topic", message)
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

	_ = r.AddSubscription("foo_topic", fooHandler)

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
