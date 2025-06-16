package kafka

import (
	"errors"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// ErrKafkaClientEmpty ...
var ErrKafkaClientEmpty = errors.New("kafka client must be set")

type kafkaClient interface {
	WriteMessages(msgs ...kafka.Message) (int, error)
}

// Kafka Service.
type Kafka struct {
	kafkaClient kafkaClient
}

// NewKafka Service constructor.
func NewKafka(kafkaClient kafkaClient) (*Kafka, error) {
	if kafkaClient == nil {
		return nil, ErrKafkaClientEmpty
	}

	return &Kafka{kafkaClient: kafkaClient}, nil
}

// ProduceInOffice Отправка сообщения о приходе сотрудника в офис / на работу.
func (k Kafka) ProduceInOffice(macAddress string) error {
	_, err := k.kafkaClient.WriteMessages(
		kafka.Message{Value: []byte(macAddress)},
	)

	if err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	return nil
}
