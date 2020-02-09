package client

import (
	"os"
	"strconv"

	"github.com/orov-io/BlackBart/response"
	api "github.com/orov-io/BlackBeard"
	"github.com/orov-io/LongBen/models"
)

const (
	portKey      = "CLIENT_TARGET_PORT"
	hostKey      = "CLIENT_TARGET_HOST"
	serviceKey   = "CLIENT_TARGET_SERVICE"
	versionKey   = "CLIENT_TARGET_VERSION"
	apiKeyKey    = "CLIENT_API_KEY"
	v1           = "v1"
	pingEndpoint = "/ping"
)

const (
	defaultHost    = "http://localhost"
	defaultVersion = "v1"
)

// Client provides functionality to easily call server endpoints
type Client struct {
	c *api.Client
}

// New returns a fresh client.
func New() *Client {
	return &Client{api.MakeNewClient()}
}

// NewWithDefaults returns a new client initialized with params from next env
// values:
// host: $CLIENT_TARGET_HOST
// port: $CLIENT_TARGET_PORT
// version: $CLIENT_TARGET_VERSION
// service: $CLIENT_TARGET_SERVICE
// api key: $CLIENT_API_KEY
func NewWithDefaults() *Client {
	host, APIKey, version, service, port := getDefaultParams()
	return New().WithBasePath(host).WithAPIKey(APIKey).
		WithVersion(version).ToService(service).WithPort(port)

}

func getDefaultParams() (host, apiKey, version, service string, port int) {
	host = getHostFromEnv()
	apiKey = getAPIKeyFromEnv()
	version = getVersionFromEnv()
	service = getServiceFromEnv()
	port = getPortFromEnv()
	return
}

func getHostFromEnv() string {
	host := os.Getenv(hostKey)
	if host == "" {
		return defaultHost
	}
	return host
}

func getAPIKeyFromEnv() string {
	return os.Getenv(apiKeyKey)
}

func getVersionFromEnv() string {
	version := os.Getenv(versionKey)
	if version == "" {
		return defaultVersion
	}
	return version
}

func getServiceFromEnv() string {
	return os.Getenv(serviceKey)
}

func getPortFromEnv() int {
	port, err := strconv.Atoi(os.Getenv(portKey))
	if err != nil {
		return 0
	}
	return port
}

// WithPort attaches desired port to underlying API Client.
func (client *Client) WithPort(port int) *Client {
	if port != 0 {
		client.c = client.c.WithPort(port)
	}
	return client
}

// WithAPIKey forces the client to send provided api key in each call to the
// server.
func (client *Client) WithAPIKey(key string) *Client {
	client.c = client.c.WithAPIKey(key)
	return client
}

// WithBasePath changes base path of the host of the underlying API Client.
func (client *Client) WithBasePath(path string) *Client {
	client.c = client.c.WithBasePath(path)
	return client
}

// WithVersion changes the version of the service of the underlying API Client.
func (client *Client) WithVersion(version string) *Client {
	client.c = client.c.WithVersion(version)
	return client
}

// ToService changes the service of the underlying API Client.
func (client *Client) ToService(service string) *Client {
	client.c = client.c.ToService(service)
	return client
}

// Ping make a call to the is_alive endpoint.
func (client *Client) Ping() (*models.Pong, error) {
	resp, err := client.c.GET(pingEndpoint, nil)
	if err != nil {
		return nil, err
	}
	pong := new(models.Pong)
	err = response.ParseTo(resp, &pong)

	return pong, err
}
