package ambient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://api.ambientweather.net"

type Device struct {
	MACAddress string
	LastData   map[string]interface{}
}

func NewClient(apiKey, appKey string, httpClient *http.Client) Client {
	return &client{
		apiKey: apiKey,
		appKey: appKey,

		httpClient: httpClient,
	}
}

type Client interface {
	Devices() ([]Device, error)
}

type client struct {
	apiKey string
	appKey string

	httpClient *http.Client
}

func (c *client) Devices() ([]Device, error) {
	var devices []Device
	if err := c.get("/v1/devices", &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

func (c *client) get(path string, out interface{}) error {
	url := fmt.Sprintf("%s%s?apiKey=%s&applicationKey=%s",
		BaseURL, path, c.apiKey, c.appKey)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return err
	}
	body := buffer.Bytes()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return json.Unmarshal(body, out)
}
