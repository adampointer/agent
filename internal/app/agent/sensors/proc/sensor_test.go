package proc

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/adampointer/agent/internal/app/agent/support"
	"github.com/prometheus/procfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPidFinder(pid int) PidFinder {
	return func(_ string) (int, error) {
		return pid, nil
	}
}

func fixturesPath(t *testing.T) string {
	return filepath.Join(support.TestDataPath(t), "fixtures", "proc")
}

func TestDefaultPidFinder(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		fsPath := fixturesPath(t)
		fs, err := procfs.NewFS(fsPath)
		require.NoError(t, err)

		f := DefaultPidFinder(fs)
		pid, err := f("vim")

		assert.NoError(t, err)
		assert.Equal(t, 26231, pid)
	})
}

func TestSensor_DoScan(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		sensor, err := NewSensorFromPath(fixturesPath(t))
		sensor.WithPidFinder(testPidFinder(26231))
		require.NoError(t, err)

		snapshot, err := sensor.DoScan(context.TODO())

		assert.NoError(t, err)
		assert.Equal(t, label, snapshot.SourceLabel)
		assert.Equal(t, &procfs.PSILine{
			Avg10:  0.1,
			Avg60:  2,
			Avg300: 3.85,
			Total:  15,
		}, snapshot.Metrics["pressure_cpu"])
		assert.Equal(t, uint(1677), snapshot.Metrics["utime"])
		assert.Equal(t, uint64(82375), snapshot.Metrics["starttime"])
		assert.Equal(t, 1.90585481e+08, snapshot.Metrics["InOctets"])
		assert.Equal(t, uint64(1418183276), snapshot.Metrics["uptime"])
		assert.Equal(t, 1038.0, snapshot.Metrics["OutRsts"])
		assert.Equal(t, &procfs.PSILine{
			Avg10:  0.1,
			Avg60:  2,
			Avg300: 3.85,
			Total:  15,
		}, snapshot.Metrics["pressure_mem"])
		assert.Equal(t, 26231, snapshot.Metrics["algod_pid"])
		assert.Equal(t, uint(44), snapshot.Metrics["stime"])
		assert.Equal(t, 7.512674e+06, snapshot.Metrics["OutOctets"])
		assert.Equal(t, 98.0, snapshot.Metrics["InErrs"])
	})
}
