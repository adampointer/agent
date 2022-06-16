package file

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutputFile(t *testing.T) {
	t.Run("file does not exist", func(t *testing.T) {
		path := newPath(t)

		f, err := outputFile(path)
		defer assert.NoError(t, f.Close())

		assert.NoError(t, err)
	})

	t.Run("file does exist", func(t *testing.T) {
		path := newPath(t)

		f1, err := os.Create(path)
		require.NoError(t, err)
		require.NoError(t, f1.Close())

		f2, err := outputFile(path)
		defer assert.NoError(t, f2.Close())

		assert.NoError(t, err)
	})
}

func TestOnEvent(t *testing.T) {
	t.Run("snapshots serialised", func(t *testing.T) {
		path := newPath(t)
		f, err := os.Create(path)

		require.NoError(t, err)

		snapshot := sensors.NewSnapshotEvent("test", map[string]interface{}{"foo": 2.2})

		assert.NoError(t, onEvent(snapshot, f))
		assert.NoError(t, f.Close())
		bs, err := os.ReadFile(path)

		expected := `"source":"test","metrics":{"foo":2.2}`
		assert.Contains(t, string(bs), expected)
	})
}

func newPath(t testing.TB) string {
	t.Helper()
	path := filepath.Join(os.TempDir(), uuid.NewString())
	t.Cleanup(func() {
		assert.NoError(t, os.Remove(path))
	})
	return path
}
