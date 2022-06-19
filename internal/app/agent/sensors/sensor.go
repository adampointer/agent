//go:generate mockgen -package=mocks -source=sensor.go -destination=../mocks/sensor.go
package sensors

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adampointer/eventbus"
)

// Sensor is a type which is able to scan some resource and return a SnapshotEvent
type Sensor interface {
	DoScan(ctx context.Context) (*SnapshotEvent, error)
}

// SnapshotEvent is a snapshot of a set of metrics at a point in time
// It implements eventbus.Event so can be sent on the eventbus
type SnapshotEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	SourceLabel string                 `json:"source"`
	Metrics     map[string]interface{} `json:"metrics"`

	topic eventbus.Topic
}

// NewSnapshotEvent creates a new SnapshotEvent from a label and a map of arbitrary metrics
// Timestamp is set by default to the time this function is called
func NewSnapshotEvent(label string, data map[string]interface{}) *SnapshotEvent {
	return &SnapshotEvent{
		SourceLabel: label,
		Timestamp:   time.Now(),
		Metrics:     data,
	}
}

// Type implements eventbus.Event
func (s *SnapshotEvent) Type() eventbus.EventType {
	return "snapshot"
}

// Topic implements eventbus.Event
func (s *SnapshotEvent) Topic() eventbus.Topic {
	return s.topic
}

// SetTopic implements eventbus.Event
func (s *SnapshotEvent) SetTopic(t eventbus.Topic) {
	s.topic = t
}

// ToJSON converts the SnapshotEvent struct to JSON
func (s *SnapshotEvent) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}
