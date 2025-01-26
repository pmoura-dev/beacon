# Beacon

Small library to interact with brokers.

## Installation

```shell
go get github.com/pmoura-dev/beacon
```

## How to use

```go
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

func fooHandler(message beacon.Message) error {
    fmt.Printf("Message received: %s\n", string(message.Payload))
    return nil
}

func main() {

    ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
    defer stop()

    mqttURL := "tcp://broker.hivemq.com:1883"

    r, err := beacon.NewRouter(
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
```
