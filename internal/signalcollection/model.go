package signalcollection

import (
	"sort"
	"sync"

	"diego.pizza/ksoc/ticker/internal/config"
	"diego.pizza/ksoc/ticker/internal/signal"
)

type SignalCollection struct {
	sync.RWMutex
	Signals []signal.Signal
}

func NewSignalCollection(cfg config.Config) SignalCollection {
	signals := cfg.Signals
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].Freq > signals[j].Freq
	})

	return SignalCollection{
		RWMutex: sync.RWMutex{},
		Signals: signals,
	}
}
