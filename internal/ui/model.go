package ui

import (
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pranavray/ikiru/internal/stats"
)

type Config struct {
	Interval   time.Duration
	Once       bool
	NoColor    bool
	Compact    bool
	ASCIIStyle string
}

type Snapshot struct {
	System  stats.SystemStats
	CPU     stats.CPUStats
	Memory  stats.MemoryStats
	Disk    stats.DiskStats
	Network stats.NetworkStats
	Battery *stats.BatteryStats
}

type Model struct {
	cfg      Config
	snap     Snapshot
	tracker  stats.NetworkTracker
	err      error
	showHelp bool
}

type tickMsg time.Time

func NewModel(cfg Config) Model {
	tracker := stats.NetworkTracker{}
	snap, err := FetchSnapshot(&tracker)
	return Model{
		cfg:     normalizeConfig(cfg),
		snap:    snap,
		tracker: tracker,
		err:     err,
	}
}

func RenderOnce(cfg Config) (string, error) {
	tracker := stats.NetworkTracker{}
	snap, _ := FetchSnapshot(&tracker)
	return RenderLayout(snap, normalizeConfig(cfg)), nil
}

func (m Model) Init() tea.Cmd {
	return tick(m.cfg.Interval)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
		case "a":
			m.cfg.ASCIIStyle = nextASCIIStyle(m.cfg.ASCIIStyle)
		case "c":
			m.cfg.Compact = !m.cfg.Compact
		case "n":
			m.cfg.NoColor = !m.cfg.NoColor
		case "+", "=", "]":
			if m.cfg.Interval < 10*time.Second {
				m.cfg.Interval += time.Second
			}
			return m, tea.ClearScreen
		case "-", "[":
			if m.cfg.Interval > time.Second {
				m.cfg.Interval -= time.Second
			}
			return m, tea.ClearScreen
		}
	case tickMsg:
		snap, err := FetchSnapshot(&m.tracker)
		m.snap = snap
		m.err = err
		return m, tick(m.cfg.Interval)
	}
	return m, nil
}

func (m Model) View() string {
	if m.showHelp {
		return RenderHelp(m.cfg) + "\n"
	}
	hint := fmt.Sprintf("\n  ? help  ·  %ds refresh", int(m.cfg.Interval.Seconds()))
	return RenderLayout(m.snap, m.cfg) + hint + "\n"
}

func FetchSnapshot(tracker *stats.NetworkTracker) (Snapshot, error) {
	var snapshot Snapshot
	var errs []error

	snapshot.System = stats.FetchSystem()

	cpuStats, err := stats.FetchCPU()
	if err != nil {
		errs = append(errs, err)
	}
	snapshot.CPU = cpuStats

	memory, err := stats.FetchMemory()
	if err != nil {
		errs = append(errs, err)
	}
	snapshot.Memory = memory

	diskStats, err := stats.FetchDisk()
	if err != nil {
		errs = append(errs, err)
	}
	snapshot.Disk = diskStats

	network, err := stats.FetchNetwork(tracker)
	if err != nil {
		errs = append(errs, err)
	}
	snapshot.Network = network

	battery, err := stats.FetchBattery()
	if err == nil && battery.Present {
		snapshot.Battery = &battery
	} else if err != nil && !errors.Is(err, stats.ErrBatteryUnavailable) {
		errs = append(errs, err)
	}

	return snapshot, errors.Join(errs...)
}

func tick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func normalizeConfig(cfg Config) Config {
	if cfg.Interval <= 0 {
		cfg.Interval = time.Second
	}
	if cfg.ASCIIStyle == "" {
		cfg.ASCIIStyle = ASCIIStylePulse
	}
	return cfg
}

func nextASCIIStyle(current string) string {
	switch current {
	case ASCIIStylePulse:
		return ASCIIStyleOS
	case ASCIIStyleOS:
		return ASCIIStyleNone
	default:
		return ASCIIStylePulse
	}
}
