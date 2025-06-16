package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAggregatorConfig_NewAggregatorConfig(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedNilObj bool
		expectedErr    error
		expectedConfig *AggregatorConfig
	}{
		"create config without OS": {
			os: func() OS {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    errors.New("OS must not be nil"),
		},
		"create config without host": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(aggregatorHostEnvParam).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", aggregatorHostEnvParam),
		},
		"create config without port": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(aggregatorHostEnvParam).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(aggregatorPortEnvParam).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", aggregatorPortEnvParam),
		},
		"create config with invalid port": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(aggregatorHostEnvParam).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(aggregatorPortEnvParam).Return("8O82").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("invalid variable %s, param must be positive integer", aggregatorPortEnvParam),
		},
		"create config": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(aggregatorHostEnvParam).Return("127.0.0.1")
				osMock.EXPECT().Getenv(aggregatorPortEnvParam).Return("8082")

				return osMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
			expectedConfig: &AggregatorConfig{
				host: "127.0.0.1",
				port: 8082,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewAggregatorConfig(test.os())
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
