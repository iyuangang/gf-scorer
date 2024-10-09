package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gf-scorer/internal/config"
	"gf-scorer/internal/database"
	"gf-scorer/internal/scorer"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	inputPath := flag.String("input", "", "Path to input file or directory")
	flag.Parse()

	if *inputPath == "" {
		log.Fatal("Input path is required")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB(db)

	s := scorer.New(db)

	inputAbs, err := filepath.Abs(*inputPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// Start metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Metrics.Port), nil))
	}()

	err = s.ProcessInput(inputAbs, cfg.Processing)
	if err != nil {
		log.Fatalf("Failed to process input: %v", err)
	}

	log.Println("Processing completed successfully")

	// Set up graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down gracefully...")
	// Perform any cleanup here
	database.CloseDB(db)
}
