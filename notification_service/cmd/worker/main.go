package main

import (
	"log"
	"notification_service/internal/config"
	"notification_service/internal/worker"
)

func main() {
	workerCfg, err := config.LoadConfig[config.WorkerConfig]("worker-config.yml")
	if err != nil {
		log.Panicf("Failed to load worker config: %s\n", err)
	}

	cfg, err := config.LoadConfig[config.Config]("config.yml")
	if err != nil {
		log.Panicf("Failed to load config: %s\n", err)
	}
	worker.WorkerRoutine(&cfg.Rabbit, &workerCfg)
}
