package utils

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"
)

type PingProbe struct {
	host     string
	interval time.Duration
}

// NewPingProbeInBackground monitors a host using ping in a go routine, once the
// packages succeeds three times, the success func is called.
func NewPingProbe(host string, intervalSeconds int) (*PingProbe, error) {
	return &PingProbe{host, time.Duration(intervalSeconds) * time.Second}, nil

}

func (p *PingProbe) Run(ctx context.Context, success func()) error {
	history := make([]bool, 3)
	online := false
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(p.interval)
		}

		c, err := net.DialTimeout("tcp", fmt.Sprintf("%s:7", p.host), time.Second)
		_ = c.Close()
		last := err == nil
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
}
