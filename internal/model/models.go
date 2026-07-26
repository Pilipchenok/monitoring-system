package model

import (
	"time"
)

type ServerMetrics struct {
	Hostname string `json:"hostname"`
	CheckTime time.Time `json:"check_time"`
	Metrics []Metric `json:"metrics"`
}

type Metric struct {
	Name string `json:"name"`
	Value float64 `json:"value"`
}

type Alert struct {
	Hostname string
	MetricName string
	Value float64
	Threshold float64
	AlertTime time.Time
}
