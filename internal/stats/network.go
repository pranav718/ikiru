package stats

import (
	"time"

	"github.com/shirou/gopsutil/v4/net"
)

type NetworkStats struct {
	BytesInPerSec  float64
	BytesOutPerSec float64
}

type NetworkTracker struct {
	lastIn  uint64
	lastOut uint64
	lastAt  time.Time
	ready   bool
}

func FetchNetwork(tracker *NetworkTracker) (NetworkStats, error) {
	counters, err := net.IOCounters(false)
	if err != nil {
		return NetworkStats{}, err
	}
	if len(counters) == 0 {
		return NetworkStats{}, nil
	}

	now := time.Now()
	currentIn := counters[0].BytesRecv
	currentOut := counters[0].BytesSent

	if tracker == nil {
		return NetworkStats{}, nil
	}

	if !tracker.ready {
		tracker.lastIn = currentIn
		tracker.lastOut = currentOut
		tracker.lastAt = now
		tracker.ready = true
		return NetworkStats{}, nil
	}

	elapsed := now.Sub(tracker.lastAt).Seconds()
	if elapsed <= 0 {
		elapsed = 1
	}

	var deltaIn, deltaOut uint64
	if currentIn >= tracker.lastIn {
		deltaIn = currentIn - tracker.lastIn
	} else {
		deltaIn = 0
	}

	if currentOut >= tracker.lastOut {
		deltaOut = currentOut - tracker.lastOut
	} else {
		deltaOut = 0
	}

	stats := NetworkStats{
		BytesInPerSec:  float64(deltaIn) / elapsed,
		BytesOutPerSec: float64(deltaOut) / elapsed,
	}

	tracker.lastIn = currentIn
	tracker.lastOut = currentOut
	tracker.lastAt = now

	return stats, nil
}
