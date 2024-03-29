package algod_metrics

import (
	"context"
	"testing"

	"github.com/adampointer/agent/internal/app/agent/support"

	"github.com/stretchr/testify/assert"
)

func TestBodyToPayload(t *testing.T) {
	t.Run("creates a payload from valid data", func(t *testing.T) {
		reader := support.GetTestDataReadClose(t, "metrics.txt")
		payload, err := bodyToPayload(reader)

		assert.NoError(t, err)
		assert.Len(t, payload, 127)
	})
}

func TestSensor_DoScan(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		srv := support.TestServer(t, path, "metrics.txt")
		sensor := NewSensorFromAddress(srv.URL)

		snapshot, err := sensor.DoScan(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, label, snapshot.SourceLabel)
		assert.Len(t, snapshot.Metrics, 127)
	})
}
