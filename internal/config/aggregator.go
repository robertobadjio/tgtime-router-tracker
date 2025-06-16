package config

import (
	"fmt"
	"net"
	"strconv"
)

const aggregatorHostEnvParam = "TGTIME_AGGREGATOR_SERVICE_HOST"
const aggregatorPortEnvParam = "TGTIME_AGGREGATOR_SERVICE_PORT"

// AggregatorConfig ...
type AggregatorConfig struct {
	host string
	port int
}

// NewAggregatorConfig ...
func NewAggregatorConfig(os OS) (*AggregatorConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("OS must not be nil")
	}

	host := os.Getenv(aggregatorHostEnvParam)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", aggregatorHostEnvParam)
	}

	portRaw := os.Getenv(aggregatorPortEnvParam)
	if len(portRaw) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", aggregatorPortEnvParam)
	}

	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid variable %s, param must be positive integer", aggregatorPortEnvParam)
	}

	return &AggregatorConfig{
		host: host,
		port: port,
	}, nil
}

// Address ...
func (ac *AggregatorConfig) Address() string {
	return net.JoinHostPort(ac.host, strconv.Itoa(ac.port))
}
