package stats

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/process"
)

type SystemStats struct {
	OS        string
	Kernel    string
	Hostname  string
	Uptime    string
	Shell     string
	Processes int
}

func FetchSystem() (SystemStats, error) {
	system := SystemStats{
		OS:        runtime.GOOS,
		Hostname:  fallbackHostname(),
		Uptime:    "unknown",
		Shell:     currentShell(),
		Processes: 0,
	}

	info, err := host.Info()
	if err == nil {
		system.OS = formatOS(info)
		system.Kernel = info.KernelVersion
		system.Hostname = info.Hostname
		system.Uptime = formatUptime(info.Uptime)
	}

	pids, err := process.Pids()
	if err != nil {
		pids = nil
	}
	system.Processes = len(pids)

	return system, nil
}

func formatOS(info *host.InfoStat) string {
	parts := []string{}
	if strings.TrimSpace(info.Platform) != "" {
		parts = append(parts, info.Platform)
	}
	if strings.TrimSpace(info.PlatformVersion) != "" {
		parts = append(parts, info.PlatformVersion)
	}
	if len(parts) == 0 {
		return runtime.GOOS
	}
	return strings.Join(parts, " ")
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	seconds %= 86400
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60

	parts := []string{}
	switch {
	case days > 0:
		parts = append(parts, formatDurationPart(days, "d"))
		if hours > 0 {
			parts = append(parts, formatDurationPart(hours, "h"))
		}
		if minutes > 0 {
			parts = append(parts, formatDurationPart(minutes, "m"))
		}
	case hours > 0:
		parts = append(parts, formatDurationPart(hours, "h"))
		if minutes > 0 {
			parts = append(parts, formatDurationPart(minutes, "m"))
		}
	default:
		parts = append(parts, formatDurationPart(minutes, "m"))
	}
	return strings.Join(parts, " ")
}

func formatDurationPart(value uint64, suffix string) string {
	return strconv.FormatUint(value, 10) + suffix
}

func currentShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" && runtime.GOOS == "windows" {
		shell = os.Getenv("COMSPEC")
	}
	if shell == "" {
		return "unknown"
	}
	return filepath.Base(shell)
}

func fallbackHostname() string {
	hostname, err := os.Hostname()
	if err != nil || strings.TrimSpace(hostname) == "" {
		return "unknown"
	}
	return hostname
}
