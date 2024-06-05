package storage

import (
	"metrics/internal/server/adapters/storage/file"
	"metrics/internal/server/adapters/storage/memory"
)

type Config struct {
	Memory *memory.Config
	File   *file.Config
}
