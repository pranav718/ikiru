package stats

import (
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUStats struct {
	Model string    `json:"model"`
	Cores []float64 `json:"cores"`
}

func FetchCPU() (CPUStats, error) {
	model := "Unknown CPU"
	if info, err := cpu.Info(); err == nil && len(info) > 0 {
		if trimmed := strings.TrimSpace(info[0].ModelName); trimmed != "" {
			model = trimmed
		}
	}

	usage, err := cpu.Percent(0, true)
	if err != nil {
		return CPUStats{Model: model}, err
	}

	return CPUStats{
		Model: model,
		Cores: usage,
	}, nil
}
