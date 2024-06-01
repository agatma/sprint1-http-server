package files

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
	"github.com/agatma/sprint1-http-server/internal/server/logger"
	"go.uber.org/zap"
)

func SaveMetricsToFile(filepath string, metrics domain.MetricValues) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create a file %w", err)
	}
	defer func(f *os.File) {
		if err = f.Close(); err != nil {
			logger.Log.Error("failed to close file: %w", zap.Error(err))
		}
	}(file)
	if err = gob.NewEncoder(file).Encode(metrics); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func LoadMetricsFromFile(filepath string) (domain.MetricValues, error) {
	var metricValues domain.MetricValues
	if _, err := os.Stat(filepath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(filepath)
		if err != nil {
			return metricValues, fmt.Errorf("failed to create file: %w", err)
		}
		err = f.Close()
		if err != nil {
			return metricValues, fmt.Errorf("failed to close file: %w", err)
		}
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		return metricValues, fmt.Errorf("failed to read file: %w", err)
	}
	if err = gob.NewDecoder(bytes.NewReader(data)).Decode(&metricValues); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed to decode file: %w", err)
		}
		return make(domain.MetricValues), nil
	}
	return metricValues, nil
}
