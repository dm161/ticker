package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	Host           string        `json:"host"`
	Port           int           `json:"port"`
	TickerInterval time.Duration `json:"ticker_interval"`
	Timeout        time.Duration `json:"timeout"`
	Signals        []Signal      `json:"signals"`
}

func ParseConfig(f io.Reader) (Config, error) {
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

type SignalCollection struct {
	sync.RWMutex
	Signals []Signal
}

func NewSignalCollection(cfg Config) SignalCollection {
	signals := cfg.Signals
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].Freq > signals[j].Freq
	})

	return SignalCollection{
		RWMutex: sync.RWMutex{},
		Signals: signals,
	}
}

type Signal struct {
	ID   uint64 `json:"id"`
	Freq uint64 `json:"frequency"`
	Msg  string `json:"message"`
}

type SignalUpdateRequest struct {
	SignalID uint64 `json:"signal_id"`
	Msg      string `json:"message"`
}

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
	cfg, err := ParseConfig(file)
	if err != nil {
		fmt.Printf("unable to parse config file %v", err)
		os.Exit(1)
	}
	signals := NewSignalCollection(cfg)
	startClock(signals, cfg.TickerInterval*time.Second, cfg.Timeout*time.Second)
	http.HandleFunc("/signal", signalHandler(signals))
	http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
}

func signalHandler(signals SignalCollection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listSignalHandler(signals)(w, r)
			return
		case http.MethodPost:
			updateSignalHandler(signals)(w, r)
			return
		default:
			fmt.Fprint(w, "method not supported")
			return
		}
	}
}

func listSignalHandler(signals SignalCollection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		signals.RLock()
		defer signals.RUnlock()
		jsonBytes, err := json.Marshal(signals.Signals)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(jsonBytes))
	}
}

func updateSignalHandler(signals SignalCollection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var updateRq SignalUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateRq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
			return
		}
		signals.Lock()
		defer signals.Unlock()
		var updated bool
		for k, v := range signals.Signals {
			if v.ID == updateRq.SignalID {
				signals.Signals[k].Msg = updateRq.Msg
				updated = true
			}
		}
		if !updated {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"status": "not found"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok"}`)
	}
}

func startClock(signals SignalCollection, tickerInterval time.Duration, timeout time.Duration) {
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
