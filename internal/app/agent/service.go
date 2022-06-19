package agent

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/adampointer/agent/internal/app/agent/settings"
	"github.com/adampointer/eventbus"
)

// Setup the application
func Setup() error {
	log.Print("setting up")

	settings.SetupEnvironment()
	if err := settings.SetupSensors(); err != nil {
		return errors.Wrap(err, "setup sensors")
	}
	settings.SetupSinks()

	return nil
}

// Start the application. This is non-blocking. Stop by canceling the context.
func Start(ctx context.Context, errC chan error) error {
	log.Print("starting")

	bus := eventbus.NewEventBus()
	start(ctx, errC, bus)

	return nil
}

func start(ctx context.Context, errC chan error, bus *eventbus.EventBus) {
	go startSinks(ctx, errC, bus)
	go startSensors(ctx, errC, settings.ScanInterval, bus)
}
