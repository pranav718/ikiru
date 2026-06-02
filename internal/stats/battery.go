package stats

import (
	"errors"
	"sync"
	"time"
)

var ErrBatteryUnavailable = errors.New("battery unavailable")

type BatteryStats struct {
	Percent  float64
	Charging bool
	State    string
	Present  bool
}

var (
	cacheMutex sync.Mutex
	cacheStats BatteryStats
	cacheErr   error
	cacheTime  time.Time
)

const cacheDuration = 30 * time.Second

func FetchBattery() (BatteryStats, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(cacheTime) < cacheDuration && cacheTime.Year() > 2000 {
		return cacheStats, cacheErr
	}

	stats, err := fetchBatteryOS()
	cacheStats = stats
	cacheErr = err
	cacheTime = time.Now()

	return stats, err
}
