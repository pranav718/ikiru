
//go:build !darwin && !linux && !windows

package stats

func fetchBatteryOS() (BatteryStats, error) {
	return BatteryStats{}, ErrBatteryUnavailable
}
