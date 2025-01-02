package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type CustomMetricError struct {
	Name   string
	Labels map[string]string
	error  error
}

func (e CustomMetricError) Unwrap() error { return e.error }

func (e CustomMetricError) Error() string {
	return fmt.Sprintf("%v, metric = %q", e.error, e.Name)
}

func GetCounter(metricName, desc string, labels map[string]string, registry *prometheus.Registry) prometheus.Counter {
	if desc == "" {
		desc = metricName
	}
	c := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: metricName,
			Help: desc,
		},
	)
	if err := registry.Register(c); err != nil {
		if resgisteredCounter, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return resgisteredCounter.ExistingCollector.(prometheus.Counter)
		}
		return nil
	}
	return c
}

func GetGauge(metricName, desc string, registry *prometheus.Registry) prometheus.Gauge {
	if desc == "" {
		desc = metricName
	}
	c := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: metricName,
			Help: desc,
		},
	)
	if err := registry.Register(c); err != nil {
		if resgisteredCounter, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return resgisteredCounter.ExistingCollector.(prometheus.Gauge)
		}
		return nil
	}
	return c
}

func GetSummary(metricName, desc string, registry *prometheus.Registry, objectives map[float64]float64) prometheus.Summary {
	if desc == "" {
		desc = metricName
	}
	s := prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       metricName,
			Help:       desc,
			Objectives: objectives,
		},
	)
	if err := registry.Register(s); err != nil {
		if resgisteredCounter, ok := err.(prometheus.AlreadyRegisteredError); ok {
			return resgisteredCounter.ExistingCollector.(prometheus.Summary)
		}
		return nil
	}
	return s
}

func UnregisterGauge(name, desc string, registry *prometheus.Registry) {
	if desc == "" {
		desc = name
	}
	c := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: desc,
		},
	)
	registry.Unregister(c)
}
