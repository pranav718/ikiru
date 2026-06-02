
//go:build linux

package stats

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func fetchBatteryOS() (BatteryStats, error) {
	matches, err := filepath.Glob("/sys/class/power_supply/BAT*")
	if err != nil || len(matches) == 0 {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	capacityRaw, err := os.ReadFile(filepath.Join(matches[0], "capacity"))
	if err != nil {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	statusRaw, err := os.ReadFile(filepath.Join(matches[0], "status"))
	if err != nil {
		statusRaw = []byte("Unknown")
	}

	percent, err := strconv.ParseFloat(strings.TrimSpace(string(capacityRaw)), 64)
	if err != nil {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	state := strings.TrimSpace(string(statusRaw))
	return BatteryStats{
		Percent:  percent,
		Charging: strings.EqualFold(state, "Charging"),
		State:    state,
		Present:  true,
	}, nil
}
