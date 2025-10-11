package config

import "time"

// BaseConfig contains fields common to both Monolith and Microservice configurations.
type BaseConfig struct {
	// Metadata
	Version   string    `yaml:"version"`
	CreatedAt time.Time `yaml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at"`

	// Project info
	ProjectName string `yaml:"project_name"`
	ProjectType string `yaml:"project_type"` // monolith, microservice

	// Cloud configuration
	CloudProvider string `yaml:"cloud_provider"` // aws, gcp, azure

}

// NewBaseConfig is a constructor for BaseConfig, setting default values.
func NewBaseConfig() BaseConfig {
	now := time.Now()
	return BaseConfig{
		Version:   "1.0",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
