package config

import "time"

// Application is application configuration.
type Application struct {
	Env string `yaml:"env" validate:"required"`

	PasswordCost int `yaml:"password_cost" validate:"required"`

	AccessTokenExp  time.Duration `yaml:"access_token_exp" validate:"required"`
	RefreshTokenExp time.Duration `yaml:"refresh_token_exp" validate:"required"`
	JwtSecret       string        `yaml:"jwt_secret" validate:"required"`

	PprofHost string `yaml:"pprof_host"`
	PprofPort int    `yaml:"pprof_port" validate:"required"`
}
