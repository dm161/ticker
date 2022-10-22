package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type SignalCollection struct {
	sync.RWMutex
	Signals []Signal
}

type Signal struct {
	ID   uint64
	Freq uint64
	Msg  string
}

type SignalUpdateRequest struct {
	SignalID uint64
	Msg      string
}

func main() {
	signals := loadSignals()
	startClock(signals, time.Second, 60*time.Second)
	http.HandleFunc("/signal", signalHandler(signals))
	http.ListenAndServe(":9876", nil)
}

func loadSignals() SignalCollection {
	return SignalCollection{
		RWMutex: sync.RWMutex{},
		Signals: []Signal{
			{ID: 1, Freq: 10, Msg: "bong"},
			{ID: 2, Freq: 3, Msg: "tock"},
			{ID: 3, Freq: 1, Msg: "tick"},
		},
	}
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
				updated := true
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
