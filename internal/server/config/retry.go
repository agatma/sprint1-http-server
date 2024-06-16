package config

import (
	"fmt"
	"metrics/internal/server/logger"
	"time"

	"github.com/avast/retry-go"
)

const addTime = 2

const Attempts uint = 3

func DelayType(n uint, _ error, config *retry.Config) time.Duration {
	switch n {
	case 0:
		return 1 * time.Second
	case 1:
		return (1 + addTime) * time.Second
	default:
		return (1 + addTime + addTime) * time.Second
	}
}

func OnRetry(n uint, err error) {
	logger.Log.Error(fmt.Sprintf(`%d %s`, n, err.Error()))
}
