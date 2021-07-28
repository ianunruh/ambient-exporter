package collect

import (
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/ianunruh/ambient-exporter/pkg/ambient"
)

func NewCollector(client ambient.Client, log *zap.Logger) *Collector {
	return &Collector{
		Client: client,
		Log:    log,

		Temp: prometheus.NewDesc("ambient_device_temp_f",
			"Temperature for a particular sensor in Fahrenheit",
			[]string{"device", "sensor"}, nil,
		),
		Humidity: prometheus.NewDesc("ambient_device_humidity",
			"Humidity percentage for a particular sensor",
			[]string{"device", "sensor"}, nil,
		),
		DewPoint: prometheus.NewDesc("ambient_device_dewpoint_f",
			"Dew point for a particular sensor in Fahrenheit",
			[]string{"device", "sensor"}, nil,
		),
		FeelsLike: prometheus.NewDesc("ambient_device_feelslike_f",
			"Feels like temperature for a particular sensor in Fahrenheit",
			[]string{"device", "sensor"}, nil,
		),
		DataAge: prometheus.NewDesc("ambient_device_data_age_seconds",
			"Age of the last data for a device in seconds",
			[]string{"device"}, nil,
		),
		Latency: prometheus.NewDesc("ambient_api_request_latency_ms",
			"Request latency for the devices API in milliseconds",
			nil, nil,
		),
	}
}

type Collector struct {
	Client ambient.Client
	Log    *zap.Logger

	Temp      *prometheus.Desc
	Humidity  *prometheus.Desc
	FeelsLike *prometheus.Desc
	DewPoint  *prometheus.Desc
	DataAge   *prometheus.Desc
	Latency   *prometheus.Desc
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Temp
	ch <- c.Humidity
	ch <- c.FeelsLike
	ch <- c.DewPoint
	ch <- c.DataAge
	ch <- c.Latency
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	devices, err := c.Client.Devices()
	latency := time.Since(start)

	if err != nil {
		c.Log.Error("Failed to call devices API",
			zap.Error(err))
		return
	}

	ch <- prometheus.MustNewConstMetric(c.Latency, prometheus.GaugeValue, float64(latency.Milliseconds()))

	for _, device := range devices {
		dataAge := time.Now().Unix() - (int64(device.LastData["dateutc"].(float64)) / 1000)
		ch <- prometheus.MustNewConstMetric(c.DataAge, prometheus.GaugeValue, float64(dataAge),
			device.MACAddress)

		for key, value := range device.LastData {
			if strings.HasPrefix(key, "temp") {
				sensor := strings.TrimSuffix(strings.TrimPrefix(key, "temp"), "f")
				ch <- prometheus.MustNewConstMetric(c.Temp, prometheus.GaugeValue, value.(float64),
					device.MACAddress, sensor)
			} else if strings.HasPrefix(key, "humidity") {
				sensor := strings.TrimPrefix(key, "humidity")
				ch <- prometheus.MustNewConstMetric(c.Humidity, prometheus.GaugeValue, value.(float64),
					device.MACAddress, sensor)
			} else if strings.HasPrefix(key, "feelsLike") {
				sensor := strings.TrimPrefix(key, "feelsLike")
				ch <- prometheus.MustNewConstMetric(c.FeelsLike, prometheus.GaugeValue, value.(float64),
					device.MACAddress, sensor)
			} else if strings.HasPrefix(key, "dewPoint") {
				sensor := strings.TrimPrefix(key, "dewPoint")
				ch <- prometheus.MustNewConstMetric(c.DewPoint, prometheus.GaugeValue, value.(float64),
					device.MACAddress, sensor)
			}
		}
	}
}
