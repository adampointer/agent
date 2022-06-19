package file

import (
	"context"
	"os"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/eventbus"
	"github.com/pkg/errors"
)

// Sink implements sinks.Sink and persists events into a log file
type Sink struct {
	path string
}

// NewSinkFromPath creates a new Sink from the path to the output file
func NewSinkFromPath(path string) *Sink {
	return &Sink{path: path}
}

// Listen for events on dataC
// This method will block until either the context is canceled or an error is returned while
// attempting to write to file. If an error is returned it will be sent on the provided errC channel
func (s *Sink) Listen(ctx context.Context, dataC chan eventbus.Event, errC chan error) {
	f, err := outputFile(s.path)
	if err != nil {
		errC <- errors.Wrap(err, "open output file")
	}

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != context.Canceled {
				errC <- ctx.Err()
				return
			}
		case evt := <-dataC:
			if err := onEvent(evt, f); err != nil {
				errC <- err
				return
			}
		}
	}
}

func outputFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func onEvent(evt eventbus.Event, f *os.File) error {
	snapshot, ok := evt.(*sensors.SnapshotEvent)
	if !ok {
		return errors.New("unable to cast event to SnapshotEvent")
	}
	return writeSnapshotToFile(f, snapshot)
}

func writeSnapshotToFile(f *os.File, snapshot *sensors.SnapshotEvent) error {
	bs, err := snapshot.ToJSON()
	if err != nil {
		return errors.Wrap(err, "serialise snapshot")
	}

	bs = append(bs, []byte("\n")...)
	_, err = f.Write(bs)

	return errors.Wrap(err, "file write")
}
