package settings

import (
	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/agent/internal/app/agent/sensors/algod_metrics"
	"github.com/adampointer/agent/internal/app/agent/sensors/algod_stats"
	"github.com/adampointer/agent/internal/app/agent/sensors/proc"
	"github.com/adampointer/agent/internal/app/agent/sinks"
	"github.com/adampointer/agent/internal/app/agent/sinks/file"
)

var (
	Sensors []sensors.Sensor
	Sinks   []sinks.Sink
)

// SetupSensors configures the standard set of Sensors used in the application
func SetupSensors() error {
	procSensor, err := proc.NewSensorFromPath(ProcFSPath)
	if err != nil {
		return err
	}
	if AlgodPidOverride > 0 {
		procSensor.WithPidFinder(func(_ string) (int, error) {
			return AlgodPidOverride, nil
		})
	}

	Sensors = []sensors.Sensor{
		algod_metrics.NewSensorFromAddress(AlgodAddress),
		algod_stats.NewSensorFromAddressAndToken(AlgodAddress, AlgodApiToken),
		procSensor,
	}
	return nil
}

// SetupSinks configures thge default set of Sinks used in the application
func SetupSinks() {
	Sinks = []sinks.Sink{
		file.NewSinkFromPath(OutputFile),
	}
}
