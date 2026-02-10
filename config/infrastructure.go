package config

// Infrastructure is a collection of infrastructure.
type Infrastructure struct {
	Database Database `yaml:"database" validate:"required"`
}
