package storage

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type MemStorage struct {
	gaugeValues   map[string]float64
	counterValues map[string]int64
	mutex         sync.RWMutex
}

var Storage = &MemStorage{
	gaugeValues:   make(map[string]float64),
	counterValues: make(map[string]int64),
}

type MetricRepository interface {
	AddGaugeValues(string, string) error
	AddCounterValues(string, string) error
}

func (m *MemStorage) String() string {
	result := make([]string, 0)
	for k, v := range m.gaugeValues {
		result = append(result, fmt.Sprintf("%s: %f", k, v))
	}
	for k, v := range m.counterValues {
		result = append(result, fmt.Sprintf("%s: %d", k, v))
	}
	return strings.Join(result, ", ")
}

func (m *MemStorage) AddGaugeValues(key, value string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	gaugeValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	m.gaugeValues[key] = gaugeValue
	return nil
}

func (m *MemStorage) AddCounterValues(key, value string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	counterValue, err := strconv.Atoi(value)
	if err != nil {
		return err
	}
	m.counterValues[key] += int64(counterValue)
	return nil
}

func (m *MemStorage) GetGaugeValues(key string) (float64, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, found := m.gaugeValues[key]
	return v, found
}

func (m *MemStorage) GetCounterValues(key string) (int64, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	v, found := m.counterValues[key]
	return v, found
}

func (m *MemStorage) GetAllGaugeValues() map[string]float64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.gaugeValues
}

func (m *MemStorage) GetAllCounterValues() map[string]int64 {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.counterValues
}
