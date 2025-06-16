package config

import (
	"fmt"
	"strconv"
	"time"
)

const (
	dbHostEnvName     = "DATABASE_HOST"
	dbPortEnvName     = "DATABASE_PORT"
	dbNameEnvName     = "DATABASE_NAME"
	dbUserEnvName     = "DATABASE_USER"
	dbPasswordEnvName = "DATABASE_PASSWORD"
	dbSSLModeEnvName  = "DATABASE_SSL_MODE"
)

const (
	dbSSLModeDisable = "disable"
	dbSSLModeRequire = "require"
)

const queryTimeoutDefault = time.Second

// DBConfig ...
type DBConfig struct {
	host         string
	port         int
	user         string
	password     string
	database     string
	sslMode      string
	queryTimeout time.Duration
}

// NewDBConfig ...
func NewDBConfig(os OS) (*DBConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("OS must not be nil")
	}

	host := os.Getenv(dbHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", dbHostEnvName)
	}

	portRaw := os.Getenv(dbPortEnvName)
	if len(portRaw) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", dbPortEnvName)
	}

	port, err := strconv.Atoi(portRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid variable %s, param must be positive integer", dbPortEnvName)
	}

	database := os.Getenv(dbNameEnvName)
	if len(database) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", dbNameEnvName)
	}

	user := os.Getenv(dbUserEnvName)
	if len(user) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", dbUserEnvName)
	}

	password := os.Getenv(dbPasswordEnvName)
	if len(password) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", dbPasswordEnvName)
	}

	SSLMode := os.Getenv(dbSSLModeEnvName)
	if len(SSLMode) == 0 {
		SSLMode = dbSSLModeDisable
	} else if SSLMode != dbSSLModeDisable && SSLMode != dbSSLModeRequire {
		return nil, fmt.Errorf("invalid DB SSL mode value")
	}

	return &DBConfig{
		host:         host,
		port:         port,
		user:         user,
		password:     password,
		database:     database,
		sslMode:      SSLMode,
		queryTimeout: queryTimeoutDefault,
	}, nil
}

// DSN ...
func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.user,
		c.password,
		c.host,
		c.port,
		c.database,
		c.sslMode,
	)
}

// QueryTimeout ...
func (c *DBConfig) QueryTimeout() time.Duration {
	return c.queryTimeout
}
