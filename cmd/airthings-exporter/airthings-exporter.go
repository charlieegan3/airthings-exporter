package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	airthings "github.com/scottlaird/airthings-exporter"
	"log"
	"net/http"
	"os"
)

type DeviceResponse struct {
	Id         string   `json:"id"`
	DeviceType string   `json:"deviceType"`
	Sensors    []string `json:"sensors"`
}

type DevicesResponse struct {
	Devices []DeviceResponse `json:"devices"`
	Offset  int              `json:"offset"`
}

// These are the values returned by my AirThings View Plus right now.
type DataValues struct {
	Battery           int     `json:"battery"`
	CO2               float32 `json:"co2"`
	Humidity          float32 `json:"humidity"`
	PM1               float32 `json:"pm1"`
	PM25              float32 `json:"pm25"`
	Pressure          float32 `json:"pressure"`
	RadonShortTermAvg float32 `json:"radonShortTermAverage"`
	Temp              float32 `json:"temp"`
	Time              int     `json:"time"`
	VOC               float32 `json:"voc"`
	RelayDeviceType   string  `json:"relayDeviceType"`
}

type DataResponse struct {
	Data DataValues `json:"data"`
}

var (
	clientID = flag.String(
		"client-id", "", "Client ID from https://dashboard.airthings.com/integrations/api-integration")
	clientSecret = flag.String(
		"client-secret", "", "Client secret from https://dashboard.airthings.com/integrations/api-integration")
)

func main() {
	flag.Parse()
	
	if (*clientID=="") || (*clientSecret=="") {
		fmt.Printf("Required --client-id and/or --client-secret parameters missing.\n")
		os.Exit(1)
	}
	
	ac := airthings.NewAPIClient(*clientID, *clientSecret)

	err := ac.AuthenticateIfNeeded()
	if err != nil {
		fmt.Printf("Authentication failed with error: %v\n", err)
		os.Exit(1)
	}

	err = ac.GetDevices()
	if err != nil {
		fmt.Printf("Failed to get device list from Airthings: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found: %d devices\n", len(ac.Devices))

	reg := prometheus.NewPedanticRegistry()
	collector := airthings.NewAirthingsCollector(ac)

	reg.MustRegister(
		collector,
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
