package config

type ServiceConfig struct {
	Environment string `mapstructure:"env"`
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
}

func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		Environment: "un-specified",
		Name:        "un-specified",
		Version:     "un-specified",
	}
}
