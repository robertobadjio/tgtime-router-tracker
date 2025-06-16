package config

import (
	"fmt"
	"net"
	"time"
)

const (
	kafkaHostEnvName = "KAFKA_HOST"
	kafkaPortEnvName = "KAFKA_PORT"
)

const connDeadlineDefault = 10 * time.Second
const inOfficeTopicDefault = "in-office"
const partitionDefault = 0

// KafkaConfig ...
type KafkaConfig struct {
	host         string
	port         string
	connDeadline time.Duration
	topic        string
	partition    int
}

// NewKafkaConfig ...
func NewKafkaConfig(os OS) (*KafkaConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("OS must not be nil")
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
		host:         host,
		port:         port,
		connDeadline: connDeadlineDefault,
		topic:        inOfficeTopicDefault,
		partition:    partitionDefault,
	}, nil
}

// Address ...
func (kc *KafkaConfig) Address() string {
	return net.JoinHostPort(kc.host, kc.port)
}

// ConnDeadline ...
func (kc *KafkaConfig) ConnDeadline() time.Duration {
	return kc.connDeadline
}

// Topic ...
func (kc *KafkaConfig) Topic() string {
	return kc.topic
}

// Partition ...
func (kc *KafkaConfig) Partition() int {
	return kc.partition
}
