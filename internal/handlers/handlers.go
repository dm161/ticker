package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"diego.pizza/ksoc/ticker/internal/signal"
)

func SignalResourceHandler(signals signal.Repo) func(w http.ResponseWriter, r *http.Request) {
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

func listSignalHandler(signals signal.Repo) func(w http.ResponseWriter, r *http.Request) {
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

func updateSignalHandler(signals signal.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var updateRq signal.SignalUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateRq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
			return
		}
		signals.Lock()
		defer signals.Unlock()
		var updated bool
		for k, v := range signals.Signals {
			if v.ID == updateRq.ID {
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
