package stats

import (
	"runtime"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStats struct {
	Path        string  `json:"path"`
	Used        uint64  `json:"used"`
	Total       uint64  `json:"total"`
	UsedPercent float64 `json:"used_percent"`
}

func FetchDisk() (DiskStats, error) {
	path := "/"
	if runtime.GOOS == "windows" {
		path = `C:\`
	}

	usage, err := disk.Usage(path)
	if err != nil {
		return DiskStats{}, err
	}

	return DiskStats{
		Path:        path,
		Used:        usage.Used,
		Total:       usage.Total,
		UsedPercent: usage.UsedPercent,
	}, nil
}
