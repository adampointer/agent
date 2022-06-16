package settings

import (
	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/agent/internal/app/agent/sensors/algod_metrics"
	"github.com/adampointer/agent/internal/app/agent/sinks"
	"github.com/adampointer/agent/internal/app/agent/sinks/file"
)

var (
	Sensors []sensors.Sensor
	Sinks   []sinks.Sink
)

func SetupSensors() {
	Sensors = []sensors.Sensor{
		algod_metrics.NewSensorFromAddress(AlgodAddress),
	}
}

func SetupSinks() {
	Sinks = []sinks.Sink{
		file.NewSink(OutputFile),
	}
}
