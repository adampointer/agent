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

type Sensor struct {
	addr string
}

func NewSensorFromAddress(addr string) *Sensor {
	return &Sensor{addr: addr}
}

func (s *Sensor) DoScan(_ context.Context) (*sensors.SnapshotEvent, error) {
	res, err := http.Get(s.addr + path)
	if err != nil {
		return nil, errors.Wrap(err, "get request")
	}
	if res.StatusCode > 299 {
		return nil, errors.Errorf("bad status code: %d (%s)", res.StatusCode, res.Status)
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
