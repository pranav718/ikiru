package stats

import (
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUStats struct {
	Model string
	Cores []float64
}

func FetchCPU() (CPUStats, error) {
	info, err := cpu.Info()
	if err != nil {
		return CPUStats{}, err
	}

	model := "Unknown CPU"
	if len(info) > 0 && strings.TrimSpace(info[0].ModelName) != "" {
		model = strings.TrimSpace(info[0].ModelName)
	}

	usage, err := cpu.Percent(120*time.Millisecond, true)
	if err != nil {
		return CPUStats{}, err
	}

	return CPUStats{
		Model: model,
		Cores: usage,
	}, nil
}
