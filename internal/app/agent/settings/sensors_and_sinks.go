package settings

import (
	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/adampointer/agent/internal/app/agent/sinks"
	"github.com/adampointer/agent/internal/app/agent/sinks/file"
)

var (
	Sensors []sensors.Sensor
	Sinks   []sinks.Sink
)

func SetupSensors() {

}

func SetupSinks() {
	Sinks = []sinks.Sink{
		file.NewSink(OutputFile),
	}
}
