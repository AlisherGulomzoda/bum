package config

// Controller is controller configuration.
type Controller struct {
	HTTP HTTPController `yaml:"http" validate:"required"`
}
