package ticker

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"diego.pizza/ksoc/ticker/internal/signal"
)

type signalRepo interface {
	List() []signal.Signal
}

type Ticker struct {
	signalRepo signalRepo
	interval   time.Duration // interval determines the base frequency at which the ticker will check for signals
	timeout    <-chan struct{}
}

func New(signalRepo signalRepo, interval time.Duration, timeout <-chan struct{}) *Ticker {
	return &Ticker{
		signalRepo: signalRepo,
		interval:   interval,
		timeout:    timeout,
	}
}

func (t *Ticker) Start(out io.Writer) <-chan struct{} {
	ticker := time.NewTicker(t.interval)
	done := make(chan struct{})
	var elapsed uint64
	go func() {
		for {
			select {
			case <-t.timeout:
				ticker.Stop()
				done <- struct{}{}
			case <-ticker.C:
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
	return done
}

func TimeoutChan(timeout time.Duration) chan struct{} {
	done := make(chan struct{})
	go func() {
		time.Sleep(timeout)
		done <- struct{}{}
	}()
	return done
}
