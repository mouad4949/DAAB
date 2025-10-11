package config

type BaseConfigApp struct {
	BaseConfig `yaml:",inline"`
	// Detection results
	Language      string   `yaml:"language"`       // nodejs, go, python, etc.
	Framework     string   `yaml:"framework"`      // express, gin, flask, etc.
	DetectedFiles []string `yaml:"detected_files"` // Files used for detection
	Port          int      `yaml:"port"`

	// Container configuration
	ContainerRegistry string `yaml:"container_registry,omitempty"` // ECR, GCR, ACR, or custom

	// Build configuration
	BuildCommand   string `yaml:"build_command,omitempty"`
	StartCommand   string `yaml:"start_command,omitempty"`
	HealthEndpoint string `yaml:"health_endpoint,omitempty"`
}

// NewBaseConfigApp is a constructor for BaseConfig, setting default values.
func NewBaseConfigApp() *BaseConfigApp {
	base := NewBaseConfig()
	return &BaseConfigApp{
		BaseConfig: base,
	}
}
