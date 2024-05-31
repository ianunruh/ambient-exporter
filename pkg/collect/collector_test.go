package collect

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/ianunruh/ambient-exporter/pkg/ambient"
)

func TestCollector(t *testing.T) {
	log := zaptest.NewLogger(t)

	devices, err := loadDevices()
	require.NoError(t, err)

	client := &mockClient{
		devices: devices,
	}

	collector := NewCollector(client, log)

	descCh := make(chan *prometheus.Desc, 6)
	collector.Describe(descCh)
	close(descCh)

	var descriptors []*prometheus.Desc
	for desc := range descCh {
		descriptors = append(descriptors, desc)
	}
	assert.Len(t, descriptors, 6)

	metricCh := make(chan prometheus.Metric, 18)
	collector.Collect(metricCh)
	close(metricCh)

	var metrics []prometheus.Metric
	for metric := range metricCh {
		fmt.Println(metric)
		metrics = append(metrics, metric)
	}
	assert.Len(t, metrics, 18)
}

type mockClient struct {
	devices []ambient.Device
}

func (c mockClient) Devices() ([]ambient.Device, error) {
	return c.devices, nil
}

func loadDevices() ([]ambient.Device, error) {
	encoded, err := os.ReadFile("testdata/devices.json")
	if err != nil {
		return nil, err
	}

	var devices []ambient.Device
	if err := json.Unmarshal(encoded, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}
