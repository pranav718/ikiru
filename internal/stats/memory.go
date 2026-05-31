package stats

import "github.com/shirou/gopsutil/v4/mem"

type MemoryStats struct {
	Used        uint64
	Total       uint64
	UsedPercent float64
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
