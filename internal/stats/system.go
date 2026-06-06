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
	OS        string `json:"os"`
	Kernel    string `json:"kernel"`
	Hostname  string `json:"hostname"`
	Uptime    string `json:"uptime"`
	Shell     string `json:"shell"`
	Processes int    `json:"processes"`
}

func FetchSystem() SystemStats {
	system := SystemStats{
		OS:        runtime.GOOS,
		Kernel:    fallbackKernel(),
		Hostname:  fallbackHostname(),
		Uptime:    "unknown",
		Shell:     currentShell(),
		Processes: 0,
	}

	info, err := host.Info()
	if err == nil {
		system.OS = formatOS(info)
		if strings.TrimSpace(info.KernelVersion) != "" {
			system.Kernel = info.KernelVersion
		}
		if strings.TrimSpace(info.Hostname) != "" {
			system.Hostname = info.Hostname
		}
		system.Uptime = formatUptime(info.Uptime)
	}

	pids, err := process.Pids()
	if err != nil {
		pids = nil
	}
	system.Processes = len(pids)

	return system
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
	seconds %= 60

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
		if seconds > 0 {
			parts = append(parts, formatDurationPart(seconds, "s"))
		}
	case hours > 0:
		parts = append(parts, formatDurationPart(hours, "h"))
		if minutes > 0 {
			parts = append(parts, formatDurationPart(minutes, "m"))
		}
		if seconds > 0 {
			parts = append(parts, formatDurationPart(seconds, "s"))
		}
	case minutes > 0:
		parts = append(parts, formatDurationPart(minutes, "m"))
		if seconds > 0 {
			parts = append(parts, formatDurationPart(seconds, "s"))
		}
	default:
		parts = append(parts, formatDurationPart(seconds, "s"))
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

func fallbackKernel() string {
	kernel, err := host.KernelVersion()
	if err != nil || strings.TrimSpace(kernel) == "" {
		return "unknown"
	}
	return kernel
}
