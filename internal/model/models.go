package model

import (
	"time"
)

type ServerMetrics struct {
	Hostname string
	CheckTime time.Time
	Metrics []Metric
}

type Metric struct {
	Name string
	Value float64
}

type Alert struct {
	Hostname string
	MetricName string
	Value float64
	Threshold float64
	AlertTime time.Time
}
