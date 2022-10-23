package ticker_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"diego.pizza/ksoc/ticker/internal/signal"
	"diego.pizza/ksoc/ticker/internal/ticker"
)

type signalRepo struct {
	signals []signal.Signal
}

func (r signalRepo) List() []signal.Signal {
	return r.signals
}

func TestTicker(t *testing.T) {
	t.Parallel()
	for name, tc := range []struct {
		interval       time.Duration
		timeout        time.Duration
		signals        []signal.Signal
		expectedOutput string
	}{
		{
			interval: time.Millisecond,
			timeout:  time.Millisecond * 550,
			signals: []signal.Signal{
				{ID: 1, Freq: 200, Msg: "tick"},
			},
			expectedOutput: "tick\ntick\n",
		},
		{
			interval: time.Millisecond,
			timeout:  time.Millisecond * 550,
			signals: []signal.Signal{
				{ID: 1, Freq: 300, Msg: "tock"},
				{ID: 2, Freq: 100, Msg: "tick"},
			},
			expectedOutput: "tick\ntick\ntock\ntick\ntick\n",
		},
		{
			interval: time.Millisecond,
			timeout:  time.Millisecond * 1050,
			signals: []signal.Signal{
				{ID: 1, Freq: 500, Msg: "beng"},
				{ID: 2, Freq: 300, Msg: "bong"},
				{ID: 3, Freq: 200, Msg: "tock"},
				{ID: 4, Freq: 100, Msg: "tick"},
			},
			expectedOutput: "tick\ntock\nbong\ntock\nbeng\nbong\ntick\ntock\nbong\nbeng\n",
		},
	} {
		t.Run(fmt.Sprintf("tc-%d", name), func(t *testing.T) {
			tckr := ticker.New(
				signalRepo{signals: tc.signals},
				tc.interval,
				ticker.TimeoutChan(tc.timeout),
			)
			buf := bytes.NewBuffer([]byte{})
			tckr.Start(buf)
			time.Sleep(tc.timeout)
			if buf.String() != tc.expectedOutput {
				t.Errorf("ticker should write '%s' to buffer, got: '%s'", tc.expectedOutput, buf.String())
			}
		})
	}
}
