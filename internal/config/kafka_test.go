package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestKafkaConfig_NewKafkaConfig(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedNilObj bool
		expectedErr    error
		expectedConfig *KafkaConfig
	}{
		"create config without OS": {
			os: func() OS {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    errors.New("os must not be nil"),
		},
		"create config without host": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(kafkaHostEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", kafkaHostEnvName),
		},
		"create config without port": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(kafkaHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(kafkaPortEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", kafkaPortEnvName),
		},
		"create config": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(kafkaHostEnvName).Return("127.0.0.1")
				osMock.EXPECT().Getenv(kafkaPortEnvName).Return("9092")

				return osMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
			expectedConfig: &KafkaConfig{
				host: "127.0.0.1",
				port: "9092",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewKafkaConfig(test.os())
			require.Equal(t, test.expectedErr, err)

			if test.expectedNilObj {
				assert.Nil(t, cfg)
			} else {
				assert.NotNil(t, cfg)
				assert.Equal(t, test.expectedConfig, cfg)
			}
		})
	}
}
