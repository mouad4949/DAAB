package initcmd

import "time"

// Config represents the DAAB configuration structure
type Config struct {
	// Metadata
	Version   string    `yaml:"version"`
	CreatedAt time.Time `yaml:"created_at"`
	UpdatedAt time.Time `yaml:"updated_at"`

	// Project info
	ProjectName string `yaml:"project_name"`
	ProjectType string `yaml:"project_type"` // monolith, microservice

	// Detection results
	Language      string   `yaml:"language"`       // nodejs, go, python, etc.
	Framework     string   `yaml:"framework"`      // express, gin, flask, etc.
	DetectedFiles []string `yaml:"detected_files"` // Files used for detection

	// Cloud configuration
	CloudProvider string `yaml:"cloud_provider"` // aws, gcp, azure
	Region        string `yaml:"region"`

	// Application configuration
	Port        int    `yaml:"port"`
	Environment string `yaml:"environment"` // production, staging, development

	// Container configuration
	ContainerRegistry string `yaml:"container_registry,omitempty"` // ECR, GCR, ACR, or custom

	// Kubernetes configuration
	Namespace string `yaml:"namespace"`

	// Build configuration
	BuildCommand   string `yaml:"build_command,omitempty"`
	StartCommand   string `yaml:"start_command,omitempty"`
	HealthEndpoint string `yaml:"health_endpoint,omitempty"`
}

func init() {
	// Set default version
}

func NewConfig() *Config {
	now := time.Now()
	return &Config{
		Version:   "1.0",
		CreatedAt: now,
		UpdatedAt: now,
	}
}
