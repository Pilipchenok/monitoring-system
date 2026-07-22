package config

import (
	"os"
	"time"
	"errors"

	"gopkg.in/yaml.v3"
)

type CollectorConfig struct {
	Port int `yaml:"port"`
	DatabaseURL string `yaml:"database_url"`
}

type AgentConfig struct {
	SendInterval time.Duration `yaml:"send_interval"`
	Hostname string `yaml:"host_name"`
	CollectorURL string `yaml:"collector_url"`
}

func LoadCollector(path string) (*CollectorConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cConf CollectorConfig
	err = yaml.Unmarshal(data, &cConf)
	if err != nil {
		return nil, err
	}

	if cConf.Port <= 0 || cConf.DatabaseURL == "" {
		return nil, errors.New("Invalid data")
	}
	return &cConf, nil
}

func LoadAgent(path string) (*AgentConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var aConf AgentConfig
	err = yaml.Unmarshal(data, &aConf)
	if err != nil {
		return nil, err
	}

	if aConf.SendInterval <= 0 || aConf.CollectorURL == "" || aConf.Hostname == "" {
		return nil, errors.New("Invalid data")
	}
	return &aConf, nil
}
