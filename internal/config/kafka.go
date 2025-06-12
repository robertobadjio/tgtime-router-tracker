package config

import (
	"fmt"
	"net"
)

const (
	kafkaHostEnvName = "KAFKA_HOST"
	kafkaPortEnvName = "KAFKA_PORT"
)

// KafkaConfig ...
type KafkaConfig struct {
	host string
	port string
}

// NewKafkaConfig ...
func NewKafkaConfig(os OS) (*KafkaConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	host := os.Getenv(kafkaHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", kafkaHostEnvName)
	}

	port := os.Getenv(kafkaPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", kafkaPortEnvName)
	}

	return &KafkaConfig{
		host: host,
		port: port,
	}, nil
}

// Address ...
func (kc *KafkaConfig) Address() string {
	return net.JoinHostPort(kc.host, kc.port)
}
