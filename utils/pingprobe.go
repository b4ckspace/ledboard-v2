package utils

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type PingProbe struct {
	host   string
	pinger *probing.Pinger
}

// NewPingProbeInBackground monitors a host using ping in a go routine, once the
// packages succeeds three times, the success func is called.
func NewPingProbe(host string, intervalSeconds int) (*PingProbe, error) {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("unable to ping: %s", err)
	}
	pinger.Interval = time.Duration(intervalSeconds) * time.Second
	return &PingProbe{host, pinger}, nil

}

func (p *PingProbe) Run(ctx context.Context, success func()) error {

	status := make(chan bool)
	defer close(status)

	go func() {
		history := make([]bool, 3)
		online := false
		for last := range status {
			history = []bool{history[1], history[2], last}

			if !online && history[0] && history[1] && history[2] {
				online = true
				slog.Info("went online", "host", p.host)
				success()
			} else if online && (!history[0] || !history[1] || !history[2]) {
				online = false
				slog.Info("went offline", "host", p.host)
			}
		}
	}()

	p.pinger.Timeout = time.Second
	p.pinger.OnRecv = func(_ *probing.Packet) {
		status <- true
	}
	p.pinger.OnRecvError = func(_ error) {
		status <- false
	}

	return p.pinger.RunWithContext(ctx)
}
