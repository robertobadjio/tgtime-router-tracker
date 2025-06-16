package config

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const intervalEnvParam = "DELAY_SECONDS"

const intervalDurationDefault = 10 * time.Second

var errorIntervalLessThanZeroOrEqual = errors.New("interval must be greater than zero")

// TrackerConfig ...
type TrackerConfig struct {
	interval time.Duration
}

// NewTrackerConfig ...
func NewTrackerConfig(os OS) (*TrackerConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("OS must not be nil")
	}

	intervalRaw := os.Getenv(intervalEnvParam)
	if len(intervalRaw) == 0 {
		return &TrackerConfig{interval: intervalDurationDefault}, nil
	}

	intervalInt, err := strconv.Atoi(intervalRaw)
	if err != nil {
		return nil, fmt.Errorf("could not parse interval from environment variable %s: %v", intervalEnvParam, err)
	}

	if intervalInt <= 0 {
		return nil, errorIntervalLessThanZeroOrEqual
	}

	return &TrackerConfig{interval: time.Duration(intervalInt) * time.Second}, nil
}

// Interval ...
func (t TrackerConfig) Interval() time.Duration {
	return t.interval
}
