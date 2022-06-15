package agent

import (
	"context"
	"time"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/eventbus"
)

const metricsTopic eventbus.Topic = "metrics"

var allSensors []sensors.Sensor

func startSensors(ctx context.Context, errC chan error, interval time.Duration, bus *eventbus.EventBus) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != context.Canceled {
				errC <- ctx.Err()
			}
			break
		case <-ticker.C:
			scan(ctx, errC, bus)
		}
	}
}

func scan(ctx context.Context, errC chan error, bus *eventbus.EventBus) {
	for _, sensor := range allSensors {
		sensor := sensor
		go func() {
			snapshot, err := sensor.DoScan(ctx)
			if err != nil {
				errC <- err
			}
			bus.Publish(metricsTopic, snapshot)
		}()
	}
}
