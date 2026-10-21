package config

import "github.com/vucongthanh92/go-base-utils/models"

// AppConfig follows the Courier service configuration contract.
type AppConfig struct {
	ServiceName     string                   `mapstructure:"serviceName"`
	Development     bool                     `mapstructure:"development"`
	Logger          *models.LoggerConfig     `mapstructure:"logger"`
	Http            *models.HttpConfig       `mapstructure:"http"`
	Database        *models.DatabaseConfig   `mapstructure:"database"`
	Tracing         *models.TracingConfig    `mapstructure:"tracing"`
	Redis           *models.RedisConfig      `mapstructure:"redis"`
	Heathcheck      *models.HeathcheckConfig `mapstructure:"heathcheck"`
	Metrics         *models.MetricsConfig    `mapstructure:"metrics"`
	CacheTTLSeconds int                      `mapstructure:"cacheTTLSeconds"`
}
