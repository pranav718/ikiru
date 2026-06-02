package stats

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var ErrBatteryUnavailable = errors.New("battery unavailable")

type BatteryStats struct {
	Percent  float64
	Charging bool
	State    string
	Present  bool
}

func FetchBattery() (BatteryStats, error) {
	switch runtime.GOOS {
	case "linux":
		return fetchLinuxBattery()
	case "darwin":
		return fetchDarwinBattery()
	default:
		return BatteryStats{}, ErrBatteryUnavailable
	}
}

func fetchLinuxBattery() (BatteryStats, error) {
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

func fetchDarwinBattery() (BatteryStats, error) {
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
