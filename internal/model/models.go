package model

import (
	"time"
)

type ServerMetrics struct {
	Name string
	CheckTime time.Time
	Metrics []Metric
}

type Metric struct {
	Name string
	Value float32
}
