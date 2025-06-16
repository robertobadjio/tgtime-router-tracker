package config

import (
	"fmt"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	routerOSVersionEnvParam = "ROUTER_OS_VERSION"
)

const (
	routerOSVersionConstraint = ">= 7.13"
	// For ROS < 7.13: /wireless/
	registrationTableWirelessSentenceDefault = "/interface/wireless/registration-table/print"
	// For ROS >= 7.13: /wifi/
	registrationTableWiFiSentenceDefault = "/interface/wifi/registration-table/print"
)

const dialTimeoutDefault = 10 * time.Second

// RouterTrackerConfig ...
type RouterTrackerConfig struct {
	registrationTableSentence string
	dialTimeout               time.Duration
}

// RegistrationTableSentence ...
func (r RouterTrackerConfig) RegistrationTableSentence() string {
	return r.registrationTableSentence
}

// DialTimeout ...
func (r RouterTrackerConfig) DialTimeout() time.Duration {
	return r.dialTimeout
}

// NewRouterTrackerConfig ...
func NewRouterTrackerConfig(os OS) (*RouterTrackerConfig, error) {
	if os == nil {
		return nil, fmt.Errorf("OS must not be nil")
	}

	routerOSVersion := registrationTableWirelessSentenceDefault
	routerOSVersionRaw := os.Getenv(routerOSVersionEnvParam)
	if len(routerOSVersionRaw) != 0 {
		wiFiInterface, err := isWiFiInterface(routerOSVersionRaw)
		if err != nil {
			return nil, fmt.Errorf("interface registration table sentence definiton error: %w", err)
		}

		if wiFiInterface {
			routerOSVersion = registrationTableWiFiSentenceDefault
		}
	}

	return &RouterTrackerConfig{
		registrationTableSentence: routerOSVersion,
		dialTimeout:               dialTimeoutDefault,
	}, nil
}

func isWiFiInterface(version string) (bool, error) {
	c, errNewConstraint := semver.NewConstraint(routerOSVersionConstraint)
	if errNewConstraint != nil {
		return false, fmt.Errorf("handle constraint not being parsable: %w", errNewConstraint)
	}

	v, errNewVersion := semver.NewVersion(version)
	if errNewVersion != nil {
		return false, fmt.Errorf("handle version not being parsable: %w", errNewVersion)
	}

	return c.Check(v), nil
}
