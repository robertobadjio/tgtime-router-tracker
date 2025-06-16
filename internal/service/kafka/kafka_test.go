package kafka

import (
	"fmt"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTimeAggregator_NewTimeAggregator(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		kafkaClient func() kafkaClient

		expectedNilObj bool
		expectedErr    error
	}{
		"create config without kafka client": {
			kafkaClient: func() kafkaClient {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    ErrKafkaClientEmpty,
		},
		"create config": {
			kafkaClient: func() kafkaClient {
				kafkaClientMock := NewMockkafkaClient(controller)
				require.NotNil(t, kafkaClientMock)

				return kafkaClientMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ta, err := NewKafka(test.kafkaClient())
			require.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, ta)
			} else {
				assert.NotNil(t, ta)
			}
		})
	}
}

func TestTimeAggregator_CreateTime(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		kafkaClient func() kafkaClient

		expectedErr error
	}{
		"produce message": {
			kafkaClient: func() kafkaClient {
				kafkaClientMock := NewMockkafkaClient(controller)
				require.NotNil(t, kafkaClientMock)

				kafkaClientMock.
					EXPECT().
					WriteMessages(kafka.Message{Value: []byte("00-1A-2B-3C-4D-5E")}).
					Return(1, nil)

				return kafkaClientMock
			},

			expectedErr: nil,
		},
		"produce message with error": {
			kafkaClient: func() kafkaClient {
				kafkaClientMock := NewMockkafkaClient(controller)
				require.NotNil(t, kafkaClientMock)

				kafkaClientMock.
					EXPECT().
					WriteMessages(kafka.Message{Value: []byte("00-1A-2B-3C-4D-5E")}).
					Return(0, fmt.Errorf("kafka internal error"))

				return kafkaClientMock
			},

			expectedErr: fmt.Errorf("failed to write messages: %w", fmt.Errorf("kafka internal error")),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			kc, errNewTimeAggregator := NewKafka(test.kafkaClient())
			require.NoError(t, errNewTimeAggregator)
			require.NotNil(t, kc)

			errCreateTime := kc.ProduceInOffice("00-1A-2B-3C-4D-5E")
			assert.Equal(t, test.expectedErr, errCreateTime)
		})
	}
}
