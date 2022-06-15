package agent

import (
	"context"
	"time"

	"github.com/adampointer/agent/internal/app/agent/settings"

	"github.com/pkg/errors"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/eventbus"
)

const metricsTopic eventbus.Topic = "metrics"

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
			allScans(ctx, errC, bus)
		}
	}
}

func allScans(ctx context.Context, errC chan error, bus *eventbus.EventBus) {
	for _, sensor := range settings.Sensors {
		sensor := sensor
		go func() {
			if err := scanAndPublish(ctx, sensor, bus); err != nil {
				errC <- err
			}
		}()
	}
}

func scanAndPublish(ctx context.Context, sensor sensors.Sensor, bus *eventbus.EventBus) error {
	snapshot, err := sensor.DoScan(ctx)
	if err != nil {
		return errors.Wrap(err, "sensor scan")
	}

	bus.Publish(metricsTopic, snapshot)
	return nil
}
