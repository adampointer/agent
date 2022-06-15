//go:generate mockgen -package=mocks -source=sink.go -destination=../mocks/sink.go
package sinks

import (
	"context"

	"github.com/adampointer/eventbus"
)

type Sink interface {
	Listen(ctx context.Context, dataC chan eventbus.Event, errC chan error)
}
