package algod_metrics

import (
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBodyToPayload(t *testing.T) {
	t.Run("creates a payload from valid data", func(t *testing.T) {
		reader := getTestDataReadClose(t)
		payload, err := bodyToPayload(reader)

		assert.NoError(t, err)
		assert.Len(t, payload, 127)
	})
}

func TestSensor_DoScan(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		srv := testServer(t)
		sensor := NewSensorFromAddress(srv.URL)

		snapshot, err := sensor.DoScan(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, label, snapshot.SourceLabel)
		assert.Len(t, snapshot.Metrics, 127)
	})
}

func getTestDataReadClose(t testing.TB) io.ReadCloser {
	t.Helper()

	data := getTestData(t)
	buf := bytes.NewBuffer(data)
	return io.NopCloser(buf)
}

func getTestData(t testing.TB) []byte {
	t.Helper()

	wd, err := os.Getwd()
	require.NoError(t, err)
	data, err := os.ReadFile(filepath.Join(wd, "..", "..", "..", "..", "..", "testdata", "metrics.txt"))
	require.NoError(t, err)
	return data
}

func testServer(t testing.TB) *httptest.Server {
	t.Helper()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/metrics", func(ginCtx *gin.Context) {
		_, err := ginCtx.Writer.Write(getTestData(t))
		require.NoError(t, err)
	})
	srv := httptest.NewServer(r)

	t.Cleanup(func() {
		srv.Close()
	})
	return srv

}
