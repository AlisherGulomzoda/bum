package bum

import (
	"embed"
)

// EmbedMigrations is a migrations embedding.
//
//go:embed migrations/*
var EmbedMigrations embed.FS
