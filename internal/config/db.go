package config

import "fmt"

const (
	dbHostEnvName     = "DATABASE_HOST"
	dbPortEnvName     = "DATABASE_PORT"
	dbNameEnvName     = "DATABASE_NAME"
	dbUserEnvName     = "DATABASE_USER"
	dbPasswordEnvName = "DATABASE_PASSWORD"
	dbSSLModeEnvName  = "DATABASE_SSL_MODE"
)

// DBConfig ...
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

// NewDBConfig ...
func NewDBConfig(os OS) (*DBConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("os must not be nil")
	}

	host := os.Getenv(dbHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", dbHostEnvName)
	}

	port := os.Getenv(dbPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", dbPortEnvName)
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
		return nil, fmt.Errorf("environment variable %s must be set", dbSSLModeEnvName)
	}

	return &DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		SSLMode:  SSLMode,
	}, nil
}

// DSN ...
func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}
