package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/adampointer/agent/internal/app/agent/mocks"
	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/agent/internal/app/agent/settings"
	"github.com/adampointer/eventbus"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSensors(t *testing.T) {
	t.Run("events are propagated over the bus", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errC := make(chan error)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		expected := sensors.NewSnapshotEvent("test", map[string]interface{}{"foo": 1.1})
		mockSensor := mocks.NewMockSensor(ctrl)
		mockSensor.
			EXPECT().
			DoScan(ctx).
			Return(expected, nil)
		settings.Sensors = []sensors.Sensor{mockSensor}

		bus := eventbus.NewEventBus()

		go startSensors(ctx, errC, time.Millisecond, bus)

		select {
		case <-time.After(time.Second):
			t.Fatal("test timeout")
		case err := <-errC:
			t.Fatal(err)
		case evt := <-bus.Subscribe(metricsTopic):
			actual, ok := evt.(*sensors.SnapshotEvent)
			assert.True(t, ok)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("context cancellation results in clean shutdown", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		errC := make(chan error)
		bus := eventbus.NewEventBus()

		go startSensors(ctx, errC, time.Millisecond, bus)
		cancel()

		select {
		case <-time.After(10 * time.Millisecond):
		case err := <-errC:
			t.Fatal(err)
		}
	})

	t.Run("sensor errors propagated", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errC := make(chan error)
		bus := eventbus.NewEventBus()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		expected := errors.New("foo")
		mockSensor := mocks.NewMockSensor(ctrl)
		mockSensor.
			EXPECT().
			DoScan(ctx).
			Return(nil, expected)
		settings.Sensors = []sensors.Sensor{mockSensor}

		go startSensors(ctx, errC, time.Millisecond, bus)

		select {
		case <-time.After(time.Second):
			t.Fatal("test timeout")
		case err := <-errC:
			assert.ErrorContains(t, err, expected.Error())
		case <-bus.Subscribe(metricsTopic):
			t.Fatal("unexpected event")
		}
	})
}
