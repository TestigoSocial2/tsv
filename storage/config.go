package storage

import (
	"os"
	"path"
)

// Config available setup options
type Config struct {
	Type       string `json:"type"`
	Path       string `json:"path"`
	Compressor string `json:"compressor"`
}

// DefaultConfig provides a sane default storage configuration
func DefaultConfig() *Config {
	return &Config{
		Type:       "bolt",
		Path:       path.Join(os.TempDir(), "bryk.db"),
		Compressor: "snappy",
	}
}
