package config

import "os"

// OS ...
type OS interface {
	Getenv(key string) string
}

type osImpl struct{}

// NewOS ...
func NewOS() OS {
	return &osImpl{}
}

// Getenv ...
func (osImpl) Getenv(key string) string {
	return os.Getenv(key)
}
