package signal

import (
	"sort"
	"sync"
)

type Repo struct {
	sync.RWMutex
	Signals []Signal
}

func NewRepo(signals []Signal) Repo {
	sort.Slice(signals, func(i, j int) bool {
		return signals[i].Freq > signals[j].Freq
	})

	return Repo{
		RWMutex: sync.RWMutex{},
		Signals: signals,
	}
}

func (r Repo) List() []Signal {
	return r.Signals
}
