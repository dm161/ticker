package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"diego.pizza/ksoc/ticker/internal/config"
	"diego.pizza/ksoc/ticker/internal/handlers"
	"diego.pizza/ksoc/ticker/internal/signalcollection"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "configpath", "", "JSON file configuration path")
	flag.Parse()
	if configPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	file, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("unable to open config file in path %v", err)
		os.Exit(1)
	}
	cfg, err := config.ParseConfig(file)
	if err != nil {
		fmt.Printf("unable to parse config file %v", err)
		os.Exit(1)
	}
	signals := signalcollection.NewSignalCollection(cfg)
	startClock(signals, cfg.TickerInterval*time.Second, cfg.Timeout*time.Second)
	http.HandleFunc("/signal", handlers.SignalResourceHandler(signals))
	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
}

func startClock(signals signalcollection.SignalCollection, tickerInterval time.Duration, timeout time.Duration) {
	ticker := time.NewTicker(tickerInterval)
	var elapsed uint64
	done := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		ticker.Stop()
		done <- struct{}{}
	}()
	go func() {
		for {
			select {
			case <-done:
				os.Exit(0)
				return
			case <-ticker.C:
				signals.RLock()
				defer signals.RUnlock()
				atomic.AddUint64(&elapsed, 1)
				for _, signal := range signals.Signals {
					if elapsed%signal.Freq == 0 {
						fmt.Println(signal.Msg)
						break
					}
				}
			}
		}
	}()
}
