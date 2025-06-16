package config

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDBConfig_NewDBConfig(t *testing.T) {
	t.Parallel()

	controller := gomock.NewController(t)

	tests := map[string]struct {
		os func() OS

		expectedNilObj bool
		expectedErr    error
		expectedConfig *DBConfig
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

				osMock.EXPECT().Getenv(dbHostEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", dbHostEnvName),
		},
		"create config without port": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(dbHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(dbPortEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", dbPortEnvName),
		},
		"create config with invalid port": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(dbHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(dbPortEnvName).Return("543O").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("invalid variable %s, param must be positive integer", dbPortEnvName),
		},
		"create config without database name": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(dbHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(dbPortEnvName).Return("5432").Times(1)
				osMock.EXPECT().Getenv(dbNameEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", dbNameEnvName),
		},
		"create config without database user": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(dbHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(dbPortEnvName).Return("5432").Times(1)
				osMock.EXPECT().Getenv(dbNameEnvName).Return("test").Times(1)
				osMock.EXPECT().Getenv(dbUserEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", dbUserEnvName),
		},
		"create config without database user password": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(dbHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(dbPortEnvName).Return("5432").Times(1)
				osMock.EXPECT().Getenv(dbNameEnvName).Return("test").Times(1)
				osMock.EXPECT().Getenv(dbUserEnvName).Return("user").Times(1)
				osMock.EXPECT().Getenv(dbPasswordEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: true,
			expectedErr:    fmt.Errorf("environment variable %s must be set", dbPasswordEnvName),
		},
		"create config without SSL mode": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(dbHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(dbPortEnvName).Return("5432").Times(1)
				osMock.EXPECT().Getenv(dbNameEnvName).Return("test").Times(1)
				osMock.EXPECT().Getenv(dbUserEnvName).Return("user").Times(1)
				osMock.EXPECT().Getenv(dbPasswordEnvName).Return("password").Times(1)
				osMock.EXPECT().Getenv(dbSSLModeEnvName).Return("").Times(1)

				return osMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
			expectedConfig: &DBConfig{
				host:     "127.0.0.1",
				port:     5432,
				user:     "user",
				password: "password",
				database: "test",
				sslMode:  dbSSLModeDisable,
			},
		},
		"create config": {
			os: func() OS {
				osMock := NewMockOS(controller)
				require.NotNil(t, osMock)

				osMock.EXPECT().Getenv(dbHostEnvName).Return("127.0.0.1").Times(1)
				osMock.EXPECT().Getenv(dbPortEnvName).Return("5432").Times(1)
				osMock.EXPECT().Getenv(dbNameEnvName).Return("test").Times(1)
				osMock.EXPECT().Getenv(dbUserEnvName).Return("user").Times(1)
				osMock.EXPECT().Getenv(dbPasswordEnvName).Return("password").Times(1)
				osMock.EXPECT().Getenv(dbSSLModeEnvName).Return(dbSSLModeRequire).Times(1)

				return osMock
			},

			expectedNilObj: false,
			expectedErr:    nil,
			expectedConfig: &DBConfig{
				host:     "127.0.0.1",
				port:     5432,
				user:     "user",
				password: "password",
				database: "test",
				sslMode:  dbSSLModeRequire,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, err := NewDBConfig(test.os())
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
