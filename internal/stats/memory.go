package stats

import "github.com/shirou/gopsutil/v4/mem"

type MemoryStats struct {
	Used        uint64  `json:"used"`
	Total       uint64  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

func FetchMemory() (MemoryStats, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return MemoryStats{}, err
	}

	return MemoryStats{
		Used:        vm.Used,
		Total:       vm.Total,
		UsedPercent: vm.UsedPercent,
	}, nil
}
