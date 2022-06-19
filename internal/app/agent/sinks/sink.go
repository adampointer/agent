//go:generate mockgen -package=mocks -source=sink.go -destination=../mocks/sink.go
package sinks

import (
	"context"

	"github.com/adampointer/eventbus"
)

// Sink is a type which can listen for SnapshotEvents on the data channel and then do something with them
type Sink interface {
	Listen(ctx context.Context, dataC chan eventbus.Event, errC chan error)
}
