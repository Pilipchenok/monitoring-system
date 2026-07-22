package config

import (
	"testing"
	"time"
)

func TestLoadCollectorConfig(t *testing.T) {
	conf, err := LoadCollector("testdata/confCollector.yaml")
	if err != nil {
		t.Fatal("Working config breaks")
	}
	if conf.Port != 8081 || conf.DatabaseURL != "postgresql://login:password@127.0.0.1:5432/database_name" {
		t.Error("Wrong config import")
	}
}

func TestEmptyCollectorConfig(t *testing.T) {
	conf, err := LoadCollector("testdata/empty_confCollector.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestUnformatCollectorConfig(t *testing.T) {
	conf, err := LoadCollector("testdata/unformat_confCollector.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestBrokenCollectorConfig(t *testing.T) {
	conf, err := LoadCollector("testdata/broken_conf.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestNonCollectorConfig(t *testing.T) {
	conf, err := LoadCollector("testdata/broken_coCollector.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of wrong config way")
	}
}

func TestLoadAgentConfig(t *testing.T) {
	conf, err := LoadAgent("testdata/confAgent.yaml")
	if err != nil {
		t.Fatal("Working config breaks")
	}
	if conf.SendInterval != time.Second * 5 ||
	conf.Hostname != "yandex.ru" ||
	conf.CollectorURL != "http://localhost:8080" {
		t.Error("Wrong config import")
	}
}

func TestEmptyAgentConfig(t *testing.T) {
	conf, err := LoadAgent("testdata/empty_confAgent.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestUnformatAgentConfig(t *testing.T) {
	conf, err := LoadAgent("testdata/unformat_confAgent.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestBrokenAgentConfig(t *testing.T) {
	conf, err := LoadAgent("testdata/broken_conf.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of empty config")
	}
}

func TestNonAgentConfig(t *testing.T) {
	conf, err := LoadAgent("testdata/broken_coAgent.yaml")
	if err == nil || conf != nil {
		t.Error("Wrong interpretation of wrong config way")
	}
}
