package utils

import (
	"log"
	"time"

	"github.com/go-ping/ping"
	"github.com/b4ckspace/ledboard-v2/config"
)

// PingProbe monitors the aliveness of a host.
type PingProbe struct {
	host               string
	alive              bool
	hostAliveCount     int
	consecutiveAnswers int
	probeInterval      time.Duration

	aliveChan chan bool
}

// NewPingProbe creates a new PingProbe instance.
func NewPingProbe(host string, cfg config.PingConfig) *PingProbe {
	return &PingProbe{
		host:               host,
		consecutiveAnswers: cfg.ConsecutiveAnswers,
		probeInterval:      time.Duration(cfg.Interval) * time.Second,
		aliveChan:          make(chan bool),
	}
}

// Start starts the ping probe.
func (p *PingProbe) Start() {
	ticker := time.NewTicker(p.probeInterval)
	defer ticker.Stop()

	for range ticker.C {
		pinger, err := ping.NewPinger(p.host)
		if err != nil {
			log.Printf("Error creating pinger: %v", err)
			continue
		}
		pinger.Count = 1 // Send only one ping packet
		pinger.Timeout = time.Second // Timeout for each ping

		err = pinger.Run()
		if err != nil {
			log.Printf("Error running pinger: %v", err)
			// Assume dead if ping fails
			p.handlePingResult(false)
			continue
		}

		stats := pinger.Statistics()
		p.handlePingResult(stats.PacketLoss == 0)
	}
}

func (p *PingProbe) handlePingResult(alive bool) {
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
func (p *PingProbe) IsAlive() bool {
	return p.alive
}

// AliveEvents returns a channel that sends true when the host becomes alive.
func (p *PingProbe) AliveEvents() <-chan bool {
	return p.aliveChan
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
