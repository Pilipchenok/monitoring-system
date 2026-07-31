package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"
	"sync"

	"monitoring-system/internal/agent"
	"monitoring-system/internal/config"
)

func main() {
	confPath := "configs/agent_config.yaml"
	cfg, err := config.LoadAgent(confPath)
	if err != nil {
		log.Fatalf("Import config error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ticker := time.NewTicker(cfg.SendInterval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				metrics, err := agent.TakeMetrics(cfg.Hostname)
				if err != nil {
					log.Printf("Taking metrics error: %v", err)
				}
				err = agent.SendMetrics(metrics, cfg.CollectorURL)
				if err != nil {
					log.Printf("Sending metrics error: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cancel()
	log.Println("Waiting for agent to stop")
	wg.Wait()

	log.Println("Agent stopped")
}
