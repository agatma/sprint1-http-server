package storage

import "strconv"

type MemStorage struct {
	gaugeValues   map[string]float64
	counterValues map[string]int64
}

var Storage = &MemStorage{
	gaugeValues:   make(map[string]float64),
	counterValues: make(map[string]int64),
}

type MetricRepository interface {
	AddGaugeValues(string, string) error
	AddCounterValues(string, string) error
}

func (m *MemStorage) AddGaugeValues(key string, value string) error {
	gaugeValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.gaugeValues[key] = gaugeValue
	return nil
}

func (m *MemStorage) AddCounterValues(key string, value string) error {
	counterValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	m.counterValues[key] += int64(counterValue)
	return nil
}
