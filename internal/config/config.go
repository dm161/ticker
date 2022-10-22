package config

import (
	"encoding/json"
	"io"
	"time"

	"diego.pizza/ksoc/ticker/internal/signal"
)

type Config struct {
	Host           string          `json:"host"`
	Port           int             `json:"port"`
	TickerInterval time.Duration   `json:"ticker_interval"`
	Timeout        time.Duration   `json:"timeout"`
	Signals        []signal.Signal `json:"signals"`
}

func ParseConfig(f io.Reader) (Config, error) {
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
