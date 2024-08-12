package migrations

import "embed"

var (
	//go:embed files/*
	Migrations embed.FS
)
