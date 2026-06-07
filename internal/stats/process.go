package stats

import (
	"sort"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPUPercent float64 `json:"cpu_percent"`
	MemPercent float32 `json:"mem_percent"`
}

var (
	procCacheMutex sync.Mutex
	procCacheData  []ProcessInfo
	procCacheTime  time.Time
)

const procCacheDuration = 5 * time.Second

func FetchTopProcesses(n int) []ProcessInfo {
	procCacheMutex.Lock()
	defer procCacheMutex.Unlock()

	if time.Since(procCacheTime) < procCacheDuration && procCacheTime.Year() > 2000 {
		return procCacheData
	}

	procs, err := process.Processes()
	if err != nil {
		return nil
	}

	grouped := make(map[string]*ProcessInfo)
	for _, p := range procs {
		name, err := p.Name()
		if err != nil || name == "" {
			continue
		}

		cpuPct, err := p.CPUPercent()
		if err != nil {
			continue
		}

		memPct, _ := p.MemoryPercent()

		if existing, ok := grouped[name]; ok {
			existing.CPUPercent += cpuPct
			existing.MemPercent += memPct
		} else {
			grouped[name] = &ProcessInfo{
				PID:        p.Pid,
				Name:       name,
				CPUPercent: cpuPct,
				MemPercent: memPct,
			}
		}
	}

	infos := make([]ProcessInfo, 0, len(grouped))
	for _, info := range grouped {
		infos = append(infos, *info)
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].CPUPercent > infos[j].CPUPercent
	})

	if len(infos) > n {
		infos = infos[:n]
	}

	procCacheData = infos
	procCacheTime = time.Now()

	return infos
}
