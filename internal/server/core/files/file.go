package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/agatma/sprint1-http-server/internal/server/core/domain"
)

func SaveMetricsToFile(file *os.File, metrics domain.MetricValues) error {
	err := file.Truncate(0)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if err = gob.NewEncoder(file).Encode(metrics); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func LoadMetricsFromFile(file *os.File) (domain.MetricValues, error) {
	var metricValues domain.MetricValues
	_, err := file.Seek(0, 0)
	if err != nil {
		return metricValues, fmt.Errorf("failed to seek file position: %w", err)
	}
	if err = gob.NewDecoder(file).Decode(&metricValues); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("failed to decode file: %w", err)
		}
		return make(domain.MetricValues), nil
	}
	return metricValues, nil
}
