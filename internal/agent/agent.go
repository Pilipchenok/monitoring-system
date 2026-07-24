package agent

import (
	"time"
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"fmt"

	"monitoring-system/internal/model"
)

var httpClient = &http.Client {
	Timeout: 5 * time.Second,
}

func TakeMetrics(hostname string) (model.ServerMetrics, error) {
	metrics := model.ServerMetrics {
    Hostname:  hostname,
    CheckTime: time.Now(),
    Metrics: []model.Metric {
			{Name: "CPU", Value: 47.7},
			{Name: "RAM", Value: 53.4},
    },
	}
	return metrics, nil
}

func SendMetrics(sm model.ServerMetrics, collectorURL string) error {
	jsonMetrics, err := json.Marshal(sm)
	if err != nil {
		log.Printf("Ошибка сериализации: %v", err)
		return err
	}

	resp, err := httpClient.Post(collectorURL + "/api/metrics", "application/json", bytes.NewBuffer(jsonMetrics))
	if err != nil {
    return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("collector returned status %d", resp.StatusCode)
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}
