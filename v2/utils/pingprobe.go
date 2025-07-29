package utils

import (
	"log/slog"
	"time"

	"github.com/prometheus-community/pro-bing"
	"github.com/b4ckspace/ledboard-v2/config"
)

// PingProbe defines the interface for monitoring host aliveness.
type PingProbe interface {
	Start()
	IsAlive() bool
	AliveEvents() <-chan bool
}

// probe implements the PingProbe interface.
type probe struct {
	host               string
	alive              bool
	hostAliveCount     int
	consecutiveAnswers int
	probeInterval      time.Duration

	aliveChan chan bool
}

// NewPingProbe creates a new PingProbe instance.
func NewPingProbe(host string, cfg config.PingConfig) PingProbe {
	return &probe{
		host:               host,
		consecutiveAnswers: cfg.ConsecutiveAnswers,
		probeInterval:      time.Duration(cfg.Interval) * time.Second,
		aliveChan:          make(chan bool),
	}
}

// Start starts the ping probe.
func (p *probe) Start() {
	ticker := time.NewTicker(p.probeInterval)
	defer ticker.Stop()

	for range ticker.C {
		pinger, err := probing.NewPinger(p.host)
		if err != nil {
			slog.Error("Error creating pinger", "error", err)
			continue
		}
		pinger.Count = 1             // Send only one ping packet
		pinger.Timeout = time.Second // Timeout for each ping

		err = pinger.Run()
		if err != nil {
			slog.Error("Error running pinger", "error", err)
			// Assume dead if ping fails
			p.handlePingResult(false)
			continue
		}

		p.handlePingResult(pinger.Statistics().PacketLoss == 0)
	}
}

func (p *probe) handlePingResult(alive bool) {
	if alive {
		if !p.alive && p.hostAliveCount >= p.consecutiveAnswers {
			p.alive = true
			p.aliveChan <- true // Emit alive event
		}
		p.hostAliveCount = min(100, p.hostAliveCount+1)
	} else {
		p.alive = false
		p.hostAliveCount = 0
		// Optionally emit dead event: p.aliveChan <- false
	}
}

// IsAlive returns true if the host is currently considered alive.
func (p *probe) IsAlive() bool {
	return p.alive
}

// AliveEvents returns a channel that sends true when the host becomes alive.
func (p *probe) AliveEvents() <-chan bool {
	return p.aliveChan
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
