package airthings

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultLabels = []string{"device_serial"}

	batteryDesc = prometheus.NewDesc(
		"airthings_device_battery_percent",
		"Current device battery level (0-100)",
		defaultLabels, nil)
	co2Desc = prometheus.NewDesc(
		"airthings_device_co2_ppm",
		"Current CO2 level in PPM",
		defaultLabels, nil)
	humDesc = prometheus.NewDesc(
		"airthings_device_humidity_perecent",
		"Current humidity level in percent",
		defaultLabels, nil)
	pm1Desc = prometheus.NewDesc(
		"airthings_device_pm1_ug_per_m3",
		"Current PM1 particulate level in micrograms per cubic meter",
		defaultLabels, nil)
	pm25Desc = prometheus.NewDesc(
		"airthings_device_pm25_ug_per_m3",
		"Current PM2.5 particulate level in micrograms per cubic meter",
		defaultLabels, nil)
	pressureDesc = prometheus.NewDesc(
		"airthings_device_pressure",
		"Current air pressure in hPa / millibars",
		defaultLabels, nil)
	radonDesc = prometheus.NewDesc(
		"airthings_device_radon",
		"24h radon average, in pCi/L",
		defaultLabels, nil)
	tempDesc = prometheus.NewDesc(
		"airthings_device_temperature_celsius",
		"Current termperature in degrees C",
		defaultLabels, nil)
	vocDesc = prometheus.NewDesc(
		"airthings_device_voc_ppb",
		"Current VOC level in parts per billion",
		defaultLabels, nil)
)

type AirthingsCollector struct {
	client *APIClient
}

func NewAirthingsCollector(client *APIClient) *AirthingsCollector {
	return &AirthingsCollector{
		client: client,
	}
}

func (c *AirthingsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *AirthingsCollector) Collect(ch chan<- prometheus.Metric) {
	err := c.client.AuthenticateIfNeeded()
	if err != nil {
		// TODO: log error
		return
	}

	err = c.client.GetDevices() // should really be cached
	if err != nil {
		// TODO: log error
		return
	}

	for _, d := range c.client.Devices {
		data, err := c.client.GetDeviceData(d)
		if err != nil {
			// TODO: log error
			return
		}
		fmt.Printf("Battery: %d, %s\n", data.Battery, d.Id)
		ch <- prometheus.MustNewConstMetric(batteryDesc, prometheus.GaugeValue, float64(data.Battery), d.Id)
		ch <- prometheus.MustNewConstMetric(co2Desc, prometheus.GaugeValue, data.CO2, d.Id)
		ch <- prometheus.MustNewConstMetric(humDesc, prometheus.GaugeValue, data.Humidity, d.Id)
		ch <- prometheus.MustNewConstMetric(pm1Desc, prometheus.GaugeValue, data.PM1, d.Id)
		ch <- prometheus.MustNewConstMetric(pm25Desc, prometheus.GaugeValue, data.PM25, d.Id)
		ch <- prometheus.MustNewConstMetric(pressureDesc, prometheus.GaugeValue, data.Pressure, d.Id)
		ch <- prometheus.MustNewConstMetric(radonDesc, prometheus.GaugeValue, data.RadonShortTermAvg, d.Id)
		ch <- prometheus.MustNewConstMetric(tempDesc, prometheus.GaugeValue, data.Temp, d.Id)
		ch <- prometheus.MustNewConstMetric(vocDesc, prometheus.GaugeValue, data.VOC, d.Id)
	}

	return
}
