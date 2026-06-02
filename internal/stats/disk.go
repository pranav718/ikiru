package stats

import (
	"runtime"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskStats struct {
	Path        string
	Used        uint64
	Total       uint64
	UsedPercent float64
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
