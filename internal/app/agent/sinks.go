package agent

import (
	"context"

	"github.com/adampointer/agent/internal/app/agent/sinks"
	"github.com/adampointer/eventbus"
)

var allSinks []sinks.Sink

func startSinks(ctx context.Context, errC chan error, bus *eventbus.EventBus) {
	dataC := bus.Subscribe(metricsTopic)

	for _, sink := range allSinks {
		go sink.Listen(ctx, dataC, errC)
	}
}
