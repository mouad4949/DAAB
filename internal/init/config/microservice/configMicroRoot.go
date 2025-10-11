package configMicroservice

import (
	config "github.com/mouad4949/DAAB/internal/init/config"
)

// Config represents the DAAB configuration structure
type ConfigMicroRoot struct {
	config.BaseConfig `yaml:",inline"` // Embed BaseConfig
	// Metadata

	DetectedMicroservices []string `yaml:"detected_files"` // Microservice
	Environment           string   `yaml:"environment"`    // production, staging, development

	// Cloud configuration
	Region string `yaml:"region"`

	Namespace string `yaml:"namespace"`
}

func NewConfigMicroRoot() *ConfigMicroRoot {
	base := config.NewBaseConfig()
	return &ConfigMicroRoot{
		BaseConfig: base,
	}
}
