package config

import (
	"time"
	//"gopkg.in/yaml.v3"
)

type CollectorConfig struct {
	Port int
	DatabaseURL string
}

type AgentConfig struct {
	SendInterval time.Duration
	HostName string
	CollectorURL string
}

func LoadCollector() {

}

func LoadAgent() {
	
}
