package config

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTrackerConfig_NewTrackerConfig(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedNilObj bool
		expectedErr    error
		expectedConfig *TrackerConfig
	}{
		"create config without OS": {
			os: func() OS {
				return nil
			},

			expectedNilObj: true,
			expectedErr:    errors.New("OS must not be nil"),
		},
		"create config with invalid interval": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(intervalEnvParam).Return("6o")

				return osMock
			},

			expectedNilObj: true,
			expectedErr: fmt.Errorf(
				"could not parse interval from environment variable %s: %v",
				intervalEnvParam,
				&strconv.NumError{Func: "Atoi", Num: "6o", Err: strconv.ErrSyntax},
			),
		},
		"create config with interval less than zero or equal": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(intervalEnvParam).Return("-1")

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    errorIntervalLessThanZeroOrEqual,
		},
		"create config without interval": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(intervalEnvParam).Return("")

				return osMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
			expectedConfig: &TrackerConfig{
				interval: intervalDurationDefault,
			},
		},
		"create config": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(intervalEnvParam).Return("60")

				return osMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
			expectedConfig: &TrackerConfig{
				interval: 60 * time.Second,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewTrackerConfig(test.os())
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
