package config_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"diego.pizza/ksoc/ticker/internal/config"
	"diego.pizza/ksoc/ticker/internal/signal"
)

func TestConfig(t *testing.T) {
	t.Parallel()
	for name, tc := range []struct {
		input  io.Reader
		output config.Config
		err    error
	}{
		{
			input:  strings.NewReader("garbage"),
			output: config.Config{},
			err:    errors.New("invalid character 'g' looking for beginning of value"),
		},
		{
			input: strings.NewReader(`{"ticker_interval": 1, "timeout": 10, "host": "localhost", "port": 1234, "signals": [{"id": 1, "frequency": 1, "message": "tick"}]}`),
			output: config.Config{
				Host:           "localhost",
				Port:           1234,
				TickerInterval: 1,
				Timeout:        10,
				Signals: []signal.Signal{
					{ID: 1, Msg: "tick", Freq: 1},
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("tc-%d", name), func(t *testing.T) {
			cfg, err := config.ParseConfig(tc.input)
			if err != nil && tc.err == nil {
				t.Errorf("expected err %v but got: %v", tc.err, err)
			}
			if err == nil && tc.err != nil {
				t.Errorf("expected err %v but got: %v", tc.err, err)
			}
			if err != nil && tc.err != nil && tc.err.Error() != err.Error() {
				t.Errorf("expected err %v but got: %v", tc.err, err)
			}
			if cfg.Host != tc.output.Host {
				t.Errorf("expected host %v but got: %v", tc.output.Host, cfg.Host)
			}
			if cfg.Port != tc.output.Port {
				t.Errorf("expected port %v but got: %v", tc.output.Port, cfg.Port)
			}
			if cfg.Timeout != tc.output.Timeout {
				t.Errorf("expected timeout %v but got: %v", tc.output.Timeout, cfg.Timeout)
			}
			if cfg.TickerInterval != tc.output.TickerInterval {
				t.Errorf("expected interval %v but got: %v", tc.output.TickerInterval, cfg.TickerInterval)
			}
			if len(cfg.Signals) != len(tc.output.Signals) {
				t.Errorf("expected signals len=%d but got: len=%d", len(tc.output.Signals), len(cfg.Signals))
			}
		})
	}
}
