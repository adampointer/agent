package agent

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/agent/internal/app/agent/settings"
	"github.com/adampointer/agent/internal/app/agent/sinks"
	"github.com/adampointer/eventbus"
	"github.com/stretchr/testify/assert"
)

type testSink struct {
	f func(event eventbus.Event)
}

func (s *testSink) Listen(_ context.Context, dataC chan eventbus.Event, _ chan error) {
	evt := <-dataC
	s.f(evt)
}

func TestSinks(t *testing.T) {
	t.Run("receives sensor events", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errC := make(chan error)
		bus := eventbus.NewEventBus()

		expected := sensors.NewSnapshotEvent("test", map[string]interface{}{"foo": 1.1})
		wg := sync.WaitGroup{}
		wg.Add(1)
		sink := &testSink{f: func(actual eventbus.Event) {
			assert.Equal(t, expected, actual)
			wg.Done()
		}}
		settings.Sinks = []sinks.Sink{sink}

		go startSinks(ctx, errC, bus)

		time.Sleep(10 * time.Millisecond)
		bus.Publish(metricsTopic, expected)

		waitWithTimeout(t, &wg, time.Second)
	})
}

func waitWithTimeout(t testing.TB, wg *sync.WaitGroup, timeout time.Duration) {
	t.Helper()

	c := make(chan struct{})

	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return
	case <-time.After(timeout):
		t.Fatal("test timeout")
	}
}
