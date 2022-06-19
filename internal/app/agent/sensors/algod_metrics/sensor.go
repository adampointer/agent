package algod_metrics

import (
	"context"
	"io"
	"net/http"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/pkg/errors"
	"github.com/prometheus/common/expfmt"
)

const (
	label = "algod_metrics"
	path  = "/metrics"
)

// Sensor implements sensors.Sensor and samples from the `/metrics` endpoint of the algod service
type Sensor struct {
	addr string
}

// NewSensorFromAddress creates a new Sensor from the algod URL
// Example: `http://localhost:8080`
func NewSensorFromAddress(addr string) *Sensor {
	return &Sensor{addr: addr}
}

// DoScan performs the scan and returns a sensors.SnapshotEvent or an error
func (s *Sensor) DoScan(_ context.Context) (*sensors.SnapshotEvent, error) {
	res, err := http.Get(s.addr + path)
	if err != nil {
		return nil, errors.Wrap(err, "get request")
	}
	if res.StatusCode > 299 {
		return nil, errors.Errorf("bad status code from metrics endpoint: %d (%s)", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	payload, err := bodyToPayload(res.Body)
	if err != nil {
		return nil, err
	}

	return sensors.NewSnapshotEvent(label, payload), nil
}

func bodyToPayload(body io.ReadCloser) (map[string]interface{}, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(body)
	if err != nil {
		return nil, errors.Wrap(err, "parse metrics text")
	}

	payload := make(map[string]interface{})
	for k, v := range mf {
		payload[k] = v.Metric
	}

	return payload, nil
}
