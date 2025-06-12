package config

import (
	"fmt"
	"net"
)

const aggregatorHostEnvParam = "TGTIME_AGGREGATOR_SERVICE_HOST"
const aggregatorPortEnvParam = "TGTIME_AGGREGATOR_SERVICE_PORT"

// AggregatorConfig ...
type AggregatorConfig struct {
	host string
	port string
}

// NewAggregatorConfig ...
func NewAggregatorConfig(os OS) (*AggregatorConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	host := os.Getenv(aggregatorHostEnvParam)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", aggregatorHostEnvParam)
	}

	port := os.Getenv(aggregatorPortEnvParam)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", aggregatorPortEnvParam)
	}

	return &AggregatorConfig{
		host: host,
		port: port,
	}, nil
}

// Address ...
func (ac *AggregatorConfig) Address() string {
	return net.JoinHostPort(ac.host, ac.port)
}
