
//go:build darwin

package stats

import (
	"os/exec"
	"strconv"
	"strings"
)

func fetchBatteryOS() (BatteryStats, error) {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	text := string(out)
	if strings.Contains(strings.ToLower(text), "no batteries") {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	percentIndex := strings.Index(text, "%")
	if percentIndex < 0 {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	start := percentIndex - 1
	for start >= 0 && text[start] >= '0' && text[start] <= '9' {
		start--
	}

	percent, err := strconv.ParseFloat(text[start+1:percentIndex], 64)
	if err != nil {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	lower := strings.ToLower(text)
	state := "Unknown"
	charging := false
	switch {
	case strings.Contains(lower, "charging"):
		state = "Charging"
		charging = true
	case strings.Contains(lower, "discharging"):
		state = "Discharging"
	case strings.Contains(lower, "charged"):
		state = "Charged"
	}

	return BatteryStats{
		Percent:  percent,
		Charging: charging,
		State:    state,
		Present:  true,
	}, nil
}
