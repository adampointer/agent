package agent

import (
	"context"
	"log"
	"time"

	"github.com/adampointer/eventbus"
)

func Setup() error {
	log.Print("setting up")
	return nil
}

func Start(ctx context.Context, errC chan error) error {
	log.Print("starting")

	bus := eventbus.NewEventBus()

	go startSinks(ctx, errC, bus)
	go startSensors(ctx, errC, time.Second, bus)

	return nil
}
