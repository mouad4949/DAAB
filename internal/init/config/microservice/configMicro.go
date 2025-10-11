package configMicroservice

import (
	config "github.com/mouad4949/DAAB/internal/init/config"
)

type ConfigMicroservice struct {
	config.BaseConfigApp `yaml:",inline"` // Embed BaseConfig
}

func NewConfigMicroservice() *ConfigMicroservice {
	baseAppConfig := config.NewBaseConfigApp()
	return &ConfigMicroservice{
		BaseConfigApp: *baseAppConfig,
	}
}
