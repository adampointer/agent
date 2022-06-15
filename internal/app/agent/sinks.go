package agent

import (
	"context"

	"github.com/adampointer/agent/internal/app/agent/settings"

	"github.com/adampointer/eventbus"
)

func startSinks(ctx context.Context, errC chan error, bus *eventbus.EventBus) {
	dataC := bus.Subscribe(metricsTopic)

	for _, sink := range settings.Sinks {
		go sink.Listen(ctx, dataC, errC)
	}
}
