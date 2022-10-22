package ticker

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"

	"diego.pizza/ksoc/ticker/internal/signal"
)

type Ticker struct {
	signalRepo signal.Repo
	interval   time.Duration // interval determines the base frequency at which the ticker will check for signals
	timeout    <-chan struct{}
}

func New(signalRepo signal.Repo, interval time.Duration, timeout <-chan struct{}) *Ticker {
	return &Ticker{
		signalRepo: signalRepo,
		interval:   interval,
		timeout:    timeout,
	}
}

func (t *Ticker) Start(out io.Writer) {
	ticker := time.NewTicker(t.interval)
	var elapsed uint64
	go func() {
		for {
			select {
			case <-t.timeout:
				os.Exit(0)
				return
			case <-ticker.C:
				t.signalRepo.RLock()
				defer t.signalRepo.RUnlock()
				atomic.AddUint64(&elapsed, 1)
				for _, signal := range t.signalRepo.List() {
					if elapsed%signal.Freq == 0 {
						fmt.Fprintf(out, "%s\n", signal.Msg)
						break
					}
				}
			}
		}
	}()
}

func TimeoutChan(timeout time.Duration) chan struct{} {
	done := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		done <- struct{}{}
	}()
	return done
}
