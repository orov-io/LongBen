package client

import (
	"os"

	"github.com/orov-io/BlackBart/response"
	api "github.com/orov-io/BlackBeard"
	"github.com/orov-io/LongBen/models"
)

const (
	portkey      = "PORT"
	serviceKey   = "SERVICE_BASE_PATH"
	v1           = "v1"
	pingEndpoint = "/ping"
)

var service = os.Getenv(serviceKey)

// Ping make a call to the is_alive endpoint.
func Ping() (*models.Pong, error) {
	client := api.MakeNewClient().WithDefaultBasePath().WithVersion(v1).
		ToService(service)
	resp, err := client.GET(pingEndpoint, nil)
	if err != nil {
		return nil, err
	}
	pong := new(models.Pong)
	err = response.ParseTo(resp, &pong)

	return pong, err
}
