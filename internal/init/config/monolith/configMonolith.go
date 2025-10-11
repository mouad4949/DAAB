package configMonolith

import (
	config "github.com/mouad4949/DAAB/internal/init/config"
)

// Config represents the DAAB configuration structure
type ConfigMonolith struct {
	config.BaseConfigApp `yaml:",inline"` // Embed BaseConfig

	// Cloud configuration
	Region string `yaml:"region"`
	// Application configuration
	Environment string `yaml:"environment"` // production, staging, development
	Namespace   string `yaml:"namespace"`
}

func init() {
	// Set default version
}

func NewConfigMonolith() *ConfigMonolith {
	// Call the constructor for BaseConfigApp to get an instance of it.
	baseAppConfig := config.NewBaseConfigApp() // Corrected: Get a *config.BaseConfigApp

	return &ConfigMonolith{
		BaseConfigApp: *baseAppConfig, // Corrected: Assign the dereferenced BaseConfigApp instance
	}
}
