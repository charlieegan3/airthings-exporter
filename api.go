package airthings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type DeviceResponse struct {
	Id         string   `json:"id"`
	DeviceType string   `json:"deviceType"`
	Sensors    []string `json:"sensors"`
}

type devicesResponse struct {
	Devices []*DeviceResponse `json:"devices"`
	Offset  int               `json:"offset"`
}

// These are the values returned by my AirThings View Plus right now.
type DataValues struct {
	Battery           int     `json:"battery"`
	CO2               float64 `json:"co2"`
	Humidity          float64 `json:"humidity"`
	PM1               float64 `json:"pm1"`
	PM25              float64 `json:"pm25"`
	Pressure          float64 `json:"pressure"`
	RadonShortTermAvg float64 `json:"radonShortTermAvg"`
	Temp              float64 `json:"temp"`
	Time              int     `json:"time"`
	VOC               float64 `json:"voc"`
	RelayDeviceType   string  `json:"relayDeviceType"`
}

type dataResponse struct {
	Data DataValues `json:"data"`
}

type APIClient struct {
	clientID, clientSecret string
	accessToken            string
	expires                time.Time
	Devices                []*DeviceResponse
}

func NewAPIClient(clientID, clientSecret string) *APIClient {
	a := &APIClient{
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	return a
}

func (a *APIClient) clearAuthentication() {
	a.expires = time.Unix(0, 0)
	a.accessToken = ""
}

func (a *APIClient) AuthenticateIfNeeded() error {
	now := time.Now()

	if a.expires.After(now) {
		// don't re-authenticate if our request hasn't expired
		return nil
	}

	resp, err := http.PostForm(
		"https://accounts-api.airthings.com/v1/token",
		url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {a.clientID},
			"client_secret": {a.clientSecret},
		})
	if err != nil {
		a.clearAuthentication()
		return err
	}

	if resp.StatusCode != 200 {
		a.clearAuthentication()
		return fmt.Errorf("Received response %q, not 200 OK.\n", resp.Status)
	}

	token := AccessTokenResponse{}
	err = json.NewDecoder(resp.Body).Decode(&token)
	defer resp.Body.Close()

	if err != nil {
		a.clearAuthentication()
		return err
	}

	a.accessToken = token.AccessToken
	a.expires = now.Add(time.Duration(token.ExpiresIn-10) * time.Second)

	return nil
}

func (a *APIClient) AddAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", a.accessToken)
}

func (a *APIClient) GetDevices() error {
	req, err := http.NewRequest("GET", "https://ext-api.airthings.com/v1/devices", nil)
	if err != nil {
		return err
	}

	a.AddAuthHeader(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("Received response %q, not 200 OK\n", resp.Status)
	}

	devices := devicesResponse{}

	err = json.NewDecoder(resp.Body).Decode(&devices)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	a.Devices = devices.Devices

	return nil
}

func (a *APIClient) GetDeviceData(device *DeviceResponse) (*DataValues, error) {
	serial := device.Id
	url := fmt.Sprintf("https://ext-api.airthings.com/v1/devices/%s/latest-samples", serial)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	a.AddAuthHeader(req)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data := dataResponse{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	return &data.Data, nil
}
