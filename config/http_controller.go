package config

// HTTPController is configuration for http controller.
type HTTPController struct {
	Address string `env:"CONTROLLER_HTTP_ADDRESS" yaml:"address"`
	Port    int    `env:"CONTROLLER_HTTP_PORT" yaml:"port" validate:"required"`
}
