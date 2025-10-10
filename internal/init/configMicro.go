package initcmd

import "time"

// Config represents the DAAB configuration structure
type ConfigMicro struct {
	// Metadata
	Version   string    `yaml:"version"`
	CreatedAt time.Time `yaml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at"`

	// Project info
	ProjectName string `yaml:"project_name"`
	ProjectType string `yaml:"project_type"` // monolith, microservice

	DetectedMicroservices []string `yaml:"detected_files"` // Files used for detection

	// Cloud configuration
	CloudProvider string `yaml:"cloud_provider"` // aws, gcp, azure

	// Build configuration
	BuildCommand   string `yaml:"build_command,omitempty"`
	StartCommand   string `yaml:"start_command,omitempty"`
	HealthEndpoint string `yaml:"health_endpoint,omitempty"`
}

func init() {
	// Set default version
}

func NewConfigMicro() *ConfigMicro {
	now := time.Now()
	return &ConfigMicro{
		Version:   "1.0",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
