package algod_stats

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/pkg/errors"
)

const (
	label        = "algod_stats"
	path         = "/v2/status"
	apiKeyHeader = "X-Algo-API-Token"
)

type Sensor struct {
	addr     string
	apiToken string
}

func NewSensorFromAddressAndToken(addr, token string) *Sensor {
	return &Sensor{
		addr:     addr,
		apiToken: token,
	}
}

func (s *Sensor) DoScan(_ context.Context) (*sensors.SnapshotEvent, error) {
	req, err := http.NewRequest(http.MethodGet, s.addr+path, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create request")
	}
	req.Header.Set(apiKeyHeader, s.apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}
	if res.StatusCode > 299 {
		return nil, errors.Errorf("bad status code from status endpoint: %d (%s)", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	payload, err := bodyToPayload(res.Body)
	if err != nil {
		return nil, err
	}

	return sensors.NewSnapshotEvent(label, payload), nil
}

func bodyToPayload(body io.ReadCloser) (map[string]interface{}, error) {
	var payload map[string]interface{}

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&payload); err != nil {
		return nil, errors.Wrap(err, "decode body")
	}

	return payload, nil
}
