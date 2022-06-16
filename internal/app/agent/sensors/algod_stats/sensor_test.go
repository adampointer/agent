package algod_stats

import (
	"context"
	"testing"

	"github.com/adampointer/agent/internal/app/agent/support"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBodyToPayload(t *testing.T) {
	t.Run("creates a payload from valid data", func(t *testing.T) {
		reader := support.GetTestDataReadClose(t, "stats.txt")
		payload, err := bodyToPayload(reader)

		assert.NoError(t, err)
		assert.Len(t, payload, 15)
	})
}

func TestSensor_DoScan(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		srv := support.TestServer(t, path, "stats.txt")
		sensor := NewSensorFromAddressAndToken(srv.URL, uuid.NewString())

		snapshot, err := sensor.DoScan(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, label, snapshot.SourceLabel)
		assert.Len(t, snapshot.Metrics, 15)
	})
}
