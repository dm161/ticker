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

func (r Repo) Update(updateRq SignalUpdateRequest) bool {
	r.Lock()
	defer r.Unlock()
	var updated bool
	for k, v := range r.Signals {
		if v.ID == updateRq.ID {
			r.Signals[k].Msg = updateRq.Msg
			updated = true
		}
	}

	return updated
}
