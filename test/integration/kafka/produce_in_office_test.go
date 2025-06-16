package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/robertobadjio/tgtime-router-tracker/internal/logger"
	kafkaInternal "github.com/robertobadjio/tgtime-router-tracker/internal/service/kafka"
	kafkaLib "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"io"
	"strings"
	"testing"
	"time"
)

const imageName = "confluentinc/confluent-local:7.5.0"

func TestProduceInOffice(t *testing.T) {
	ctx := context.Background()
	kafkaContainer, err := kafka.Run(
		ctx,
		imageName,
		kafka.WithClusterID("test-cluster"),
	)
	require.NoError(t, err)
	defer func() {
		errTerminateKafka := testcontainers.TerminateContainer(kafkaContainer)
		require.NoError(t, errTerminateKafka)
	}()

	state, errGetState := kafkaContainer.State(ctx)
	require.NoError(t, errGetState)
	t.Log("Kafka cluster name:", kafkaContainer.ClusterID)
	t.Log("Kafka cluster state running:", state.Running)

	kafkaBrokers, errGetBrokers := kafkaContainer.Brokers(ctx)
	require.NoError(t, errGetBrokers, "Error getting kafka brokers")
	require.NotEmpty(t, kafkaBrokers, "Empty kafka brokers")
	t.Log("Kafka brokers:", strings.Join(kafkaBrokers, ", "))

	time.Sleep(time.Second)

	conn, errDialLeader := kafkaLib.DialLeader(
		ctx,
		"tcp",
		kafkaBrokers[0],
		"test-topic",
		0,
	)
	if errDialLeader != nil {
		logger.Fatal(
			"component", "di",
			"during", "dial leader kafka",
			"err", errDialLeader.Error(),
		)
	}
	defer func() { _ = conn.Close() }()

	errSetWriteDeadline := conn.SetWriteDeadline(time.Now().Add(time.Second))
	if errSetWriteDeadline != nil {
		logger.Fatal(
			"component", "di",
			"during", "kafka set deadline",
			"err", errSetWriteDeadline.Error(),
		)
	}

	kafkaClient, errNewKafka := kafkaInternal.NewKafka(conn)
	require.NoError(t, errNewKafka)

	errProduceInOffice := kafkaClient.ProduceInOffice("00:1A:2B:3C:4D:5E")
	require.NoError(t, errProduceInOffice)

	gotMacAddress, errConsumeInOffice := ConsumeInOffice(ctx, kafkaBrokers)
	require.NoError(t, errConsumeInOffice)

	assert.Equal(t, "00:1A:2B:3C:4D:5E", gotMacAddress)
}

func ConsumeInOffice(ctx context.Context, brokers []string) (string, error) {
	r := buildReader("in-office", brokers)
	defer func() {
		if err := r.Close(); err != nil {
			logger.Warn(
				"component", "kafka",
				"during", "consume in office",
				"desc", "failed to close reader",
				"error", err.Error(),
			)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("exit from consumer in office: %w", ctx.Err())
		default:
		}

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("exit from consumer in office: %w", ctx.Err())
		default:
			m, err := r.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, io.EOF) {
					break // TODO: ?
				}
				return "", fmt.Errorf("reading message: %w", err)
			}

			return string(m.Value), nil
		}
	}
}

func buildReader(topicName string, brokers []string) *kafkaLib.Reader {
	return kafkaLib.NewReader(kafkaLib.ReaderConfig{
		Brokers:   brokers,
		Topic:     "in-office",
		Partition: 0,
		GroupID:   "",
		MaxBytes:  10e3,
	})
}
