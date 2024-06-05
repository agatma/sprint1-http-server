package storage

import "metrics/internal/agent/adapters/storage/memory"

type Config struct {
	Memory *memory.Config
}
