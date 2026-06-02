
//go:build windows

package stats

import (
	"syscall"
	"unsafe"
)

type systemPowerStatus struct {
	ACLineStatus        byte
	BatteryFlag         byte
	BatteryLifePercent  byte
	SystemStatusFlag    byte
	BatteryLifeTime     uint32
	BatteryFullLifeTime uint32
}

func fetchBatteryOS() (BatteryStats, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getSystemPowerStatus := kernel32.NewProc("GetSystemPowerStatus")

	var status systemPowerStatus
	r1, _, err := getSystemPowerStatus.Call(uintptr(unsafe.Pointer(&status)))
	if r1 == 0 {
		return BatteryStats{}, err
	}

	if status.BatteryFlag == 128 || status.BatteryLifePercent == 255 {
		return BatteryStats{}, ErrBatteryUnavailable
	}

	percent := float64(status.BatteryLifePercent)
	isCharging := (status.BatteryFlag & 8) != 0

	state := "Discharging"
	if isCharging {
		state = "Charging"
	} else if status.ACLineStatus == 1 {
		state = "Plugged In"
		if percent == 100 {
			state = "Charged"
		}
	}

	return BatteryStats{
		Percent:  percent,
		Charging: isCharging || (status.ACLineStatus == 1 && percent < 100),
		State:    state,
		Present:  true,
	}, nil
}
