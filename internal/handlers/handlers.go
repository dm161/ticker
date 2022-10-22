package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"diego.pizza/ksoc/ticker/internal/signal"
)

func SignalResourceHandler(signalRepo signal.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listSignalHandler(signalRepo)(w, r)
			return
		case http.MethodPost:
			updateSignalHandler(signalRepo)(w, r)
			return
		default:
			fmt.Fprint(w, "method not supported")
			return
		}
	}
}

func listSignalHandler(signalRepo signal.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		signalRepo.RLock()
		defer signalRepo.RUnlock()
		jsonBytes, err := json.Marshal(signalRepo.List())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(jsonBytes))
	}
}

func updateSignalHandler(signalRepo signal.Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var updateRq signal.SignalUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updateRq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error": "%v"}`, err.Error())
			return
		}
		if !signalRepo.Update(updateRq) {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"status": "not found"}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status": "ok"}`)
	}
}
