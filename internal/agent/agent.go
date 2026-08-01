package agent

import (
	"time"
	"encoding/json"
	"log"
	"net/http"
	"bytes"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/net"

	"monitoring-system/internal/model"
)

var httpClient = &http.Client {
	Timeout: 5 * time.Second,
}

func TakeMetrics(hostname string) (model.ServerMetrics, error) {
	var collectedMetrics []model.Metric
	vMem, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("RAM scanning error: %v", err)
	} else {
		collectedMetrics = append(collectedMetrics, model.Metric{Name: "RAM", Value: vMem.UsedPercent})
	}
	cpuPercents, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercents) == 0 {
		log.Printf("CPU scanning error: %v", err)
	} else {
		collectedMetrics = append(collectedMetrics, model.Metric{Name: "CPU", Value: cpuPercents[0]})
	}
	usage, err := disk.Usage("/")
	if err != nil {
		log.Printf("Disk scanning error: %v", err)
	} else {
		collectedMetrics = append(collectedMetrics, model.Metric{Name: "DISK", Value: usage.UsedPercent})
	}
	connections, err := net.Connections("all")
	if err != nil {
		log.Printf("Net scanning error: %v", err)
	} else {
		collectedMetrics = append(collectedMetrics, model.Metric{Name: "ActiveConnections", Value: float64(len(connections))})
	}
	metrics := model.ServerMetrics {
    Hostname:  hostname,
    CheckTime: time.Now(),
    Metrics: collectedMetrics,
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
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("collector returned status %d", resp.StatusCode)
		return fmt.Errorf("bad status: %d", resp.StatusCode)
	}
	return nil
}
