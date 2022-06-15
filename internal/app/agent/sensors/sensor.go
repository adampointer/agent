package sensors

import (
	"context"
	"encoding/json"
	"time"

	"github.com/adampointer/eventbus"
)

type SnapshotEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	SourceLabel string                 `json:"source"`
	Metrics     map[string]interface{} `json:"metrics"`

	topic eventbus.Topic
}

func NewSnapshotEvent(label string, data map[string]interface{}) *SnapshotEvent {
	return &SnapshotEvent{
		SourceLabel: label,
		Timestamp:   time.Now(),
		Metrics:     data,
	}
}

func (s *SnapshotEvent) Type() eventbus.EventType {
	return "snapshot"
}

func (s *SnapshotEvent) Topic() eventbus.Topic {
	return s.topic
}

func (s *SnapshotEvent) SetTopic(t eventbus.Topic) {
	s.topic = t
}

func (s *SnapshotEvent) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}

type Sensor interface {
	DoScan(ctx context.Context) (*SnapshotEvent, error)
}
