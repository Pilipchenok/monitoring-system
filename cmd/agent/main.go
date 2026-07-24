package main

import (
	"log"
	"time"

	"monitoring-system/internal/agent"
	"monitoring-system/internal/config"
)

func main() {
	confPath := "configs/agent_config.yaml"
	cfg, err := config.LoadAgent(confPath)
	if err != nil {
		log.Fatalf("Import config error: %v", err)
	}

	go func() {
		for {
			metrics, err := agent.TakeMetrics(cfg.Hostname)
			if err != nil {
				log.Printf("Taking metrics error: %v", err)
			}
			err = agent.SendMetrics(metrics, cfg.CollectorURL)
			if err != nil {
				log.Printf("Sending metrics error: %v", err)
			}
			time.Sleep(cfg.SendInterval)
		}
	}()


}
