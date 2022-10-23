package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"diego.pizza/ksoc/ticker/internal/config"
	"diego.pizza/ksoc/ticker/internal/handlers"
	"diego.pizza/ksoc/ticker/internal/signal"
	"diego.pizza/ksoc/ticker/internal/ticker"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "configpath", "", "JSON file configuration path")
	flag.Parse()
	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(configPath); err != nil {
		log.Fatal(err)
	}
}

func run(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("unable to open config file in path %v", err)
	}
	cfg, err := config.ParseConfig(file)
	if err != nil {
		return fmt.Errorf("unable to parse config file %v", err)
	}
	signals := signal.NewRepo(cfg.Signals)
	tc := ticker.New(
		signals,
		cfg.TickerInterval*time.Second,
		ticker.TimeoutChan(cfg.Timeout*time.Second),
	)
	done := tc.Start(os.Stdout)
	go func() {
		<-done
		os.Exit(0)
	}()
	http.HandleFunc("/signal", handlers.SignalResourceHandler(signals))
	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
	return nil
}
