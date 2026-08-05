package main

import (
	"log"
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"context"
	"time"

	"monitoring-system/internal/collector/database"
	"monitoring-system/internal/collector/handler"
	"monitoring-system/internal/collector/repository"
	"monitoring-system/internal/collector/service"
	"monitoring-system/internal/config"
)

func main() {
	confPath := "configs/collector_config.yaml"
	cfg, err := config.LoadCollector(confPath)
	if err != nil {
		log.Fatalf("Import config error: %v", err)
	}

	sqlDB, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Connect to database error: %v", err)
	}
	defer sqlDB.Close()

	rep := repository.NewRepository(sqlDB)
	serv := service.NewService(rep)
	myHandler := handler.NewHandler(serv)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/metrics", myHandler.SaveMetricsHandler)
	mux.HandleFunc("GET /api/metrics/latest", myHandler.GetMetricsHandler)
	mux.HandleFunc("GET /api/hosts", myHandler.GetHostsHandler)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/index.html")
	})

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ticker := time.NewTicker(3 * time.Hour)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				threshold := time.Now().AddDate(0, 0, -7)
				err := serv.CleanOldMetrics(ctx, threshold)
				if err != nil {
					log.Printf("Cleaning metrics error: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		log.Printf("Collector started on :%d", cfg.Port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	errSD := server.Shutdown(shutdownCtx)
	if errSD != nil {
		log.Fatalf("Stop server error: %v", errSD)
	}

	log.Println("Server stopped")
}
