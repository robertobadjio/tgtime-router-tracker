package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// Kafka Клиент для кафки
type Kafka struct {
	address string
}

// NewKafka Конструктор клиента
func NewKafka(address string) *Kafka {
	return &Kafka{address: address}
}

// ProduceInOffice Отправка сообщения о приходе сотрудника в офис / на работу
func (k Kafka) ProduceInOffice(ctx context.Context, macAddress string) error {
	conn, err := kafka.DialLeader(
		ctx,
		"tcp",
		k.address,
		inOfficeTopic,
		partition,
	)
	if err != nil {
		return fmt.Errorf("failed to dial leader: %w", err)
	}

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return fmt.Errorf("failed to set deadline: %w", err)
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(macAddress)},
	)
	if err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	if err = conn.Close(); err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}
