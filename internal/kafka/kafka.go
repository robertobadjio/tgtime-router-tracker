package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	address string
}

func NewKafka(address string) *Kafka {
	return &Kafka{address: address}
}

func (k Kafka) Produce(ctx context.Context, m InOfficeMessage, topic string) error {
	conn, err := kafka.DialLeader(
		ctx,
		"tcp",
		k.address,
		topic,
		partition,
	)
	if err != nil {
		return err
	}

	/*err = conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	if err != nil {
		return err
	}*/

	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(m.MacAddress)},
	)
	if err != nil {
		return err
	}

	if err = conn.Close(); err != nil {
		return err
	}

	return nil
}
