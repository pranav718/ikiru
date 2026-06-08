package ui

import (
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pranavray/ikiru/internal/stats"
	"github.com/pranavray/ikiru/internal/theme"
)

func RenderLayout(snapshot, prevSnap Snapshot, cfg Config, width, height int, hint string) string {
	coresPerRow := 2
	barWidth := 10
	if width >= 120 {
		coresPerRow = 4
		barWidth = 14
	} else if width >= 100 {
		coresPerRow = 3
		barWidth = 12
	}

	statsView := renderStats(snapshot, prevSnap, cfg, coresPerRow, barWidth, height)
	if hint != "" {
		statsView += "\n" + hint
	}
	if cfg.Compact || cfg.ASCIIStyle == ASCIIStyleNone {
		return statsView
	}

	art := ASCIIArt(cfg.ASCIIStyle)
	artHeight := lipgloss.Height(art)
	statsHeight := lipgloss.Height(statsView)
	if artHeight < statsHeight {
		art = strings.Repeat("\n", (statsHeight-artHeight)/2) + art
	}

	left := lipgloss.NewStyle().
		PaddingLeft(2).
		PaddingRight(1).
		Render(style(cfg, theme.SpringBlue).Render(art))
	right := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(color(cfg, theme.SumiInk4)).
		PaddingLeft(1).
		Render(statsView)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}


func renderStats(snapshot, prevSnap Snapshot, cfg Config, coresPerRow, barWidth, height int) string {
	memTrend := trend(snapshot.Memory.UsedPercent, prevSnap.Memory.UsedPercent, cfg)
	diskTrend := trend(snapshot.Disk.UsedPercent, prevSnap.Disk.UsedPercent, cfg)
	netInTrend := trend(snapshot.Network.BytesInPerSec, prevSnap.Network.BytesInPerSec, cfg)
	netOutTrend := trend(snapshot.Network.BytesOutPerSec, prevSnap.Network.BytesOutPerSec, cfg)

	coreRows := (len(snapshot.CPU.Cores) + coresPerRow - 1) / coresPerRow

	lines := []string{
		row("OS", snapshot.System.OS+"  "+snapshot.System.Kernel, cfg),
		row("Host", snapshot.System.Hostname, cfg),
		row("Uptime", snapshot.System.Uptime, cfg),
		row("Shell", snapshot.System.Shell, cfg),
		row("CPU", snapshot.CPU.Model, cfg),
		row("Cores", renderCoreGrid(snapshot.CPU.Cores, prevSnap.CPU.Cores, cfg, coresPerRow), cfg),
		row("Memory", fmt.Sprintf("%s / %s %s %s", bytes(snapshot.Memory.Used), bytes(snapshot.Memory.Total), bar(snapshot.Memory.UsedPercent, barWidth, cfg), memTrend), cfg),
		row("Disk", fmt.Sprintf("%s / %s %s %s", bytes(snapshot.Disk.Used), bytes(snapshot.Disk.Total), bar(snapshot.Disk.UsedPercent, barWidth, cfg), diskTrend), cfg),
		row("Network", fmt.Sprintf("↑ %s/s %s  ↓ %s/s %s", bytes(uint64(snapshot.Network.BytesOutPerSec)), netOutTrend, bytes(uint64(snapshot.Network.BytesInPerSec)), netInTrend), cfg),
		row("Processes", fmt.Sprintf("%d", snapshot.System.Processes), cfg),
	}

	if snapshot.Battery != nil && snapshot.Battery.Present {
		lines = append(lines, row("Battery", fmt.Sprintf("%.0f%% %s %s", snapshot.Battery.Percent, snapshot.Battery.State, batteryBar(snapshot.Battery.Percent, 8, cfg)), cfg))
	}

	// base line count: 10 stats + coreRows-1 extra lines from grid + 1 battery (if present)
	// + 1 hint line always added by caller
	baseLines := len(lines) + coreRows - 1 + 1

	// calculate how many optional sections fit
	// processes section: 1 blank + len(procs) lines
	// palette section: 1 blank + 1 palette line
	procsLines := 0
	if len(snapshot.TopProcs) > 0 {
		procsLines = 1 + len(snapshot.TopProcs)
	}
	paletteLines := 2

	showProcs := len(snapshot.TopProcs) > 0
	showPalette := true

	if height > 0 {
		available := height - baseLines
		if available < procsLines+paletteLines {
			// not enough room for both, drop palette first
			showPalette = false
		}
		if available < procsLines {
			// still not enough, drop processes too
			showProcs = false
		}
	}

	if showProcs {
		lines = append(lines, "", renderProcesses(snapshot.TopProcs, cfg))
	}

	if showPalette {
		lines = append(lines, "", palette(cfg))
	}

	return strings.Join(lines, "\n")
}

func row(label string, value string, cfg Config) string {
	labelText := fmt.Sprintf("%-10s", label)
	return style(cfg, theme.OniViolet).Bold(true).Render(labelText) + style(cfg, theme.FujiWhite).Render(value)
}

func renderCoreGrid(cores, prevCores []float64, cfg Config, coresPerRow int) string {
	if len(cores) == 0 {
		return "unavailable"
	}

	items := make([]string, 0, len(cores))
	for i, usage := range cores {
		var prevUsage float64
		if i < len(prevCores) {
			prevUsage = prevCores[i]
		}
		t := trend(usage, prevUsage, cfg)
		items = append(items, fmt.Sprintf("c%d %s%3.0f%%%s", i, bar(usage, 5, cfg), usage, t))
	}

	rows := []string{}
	for i := 0; i < len(items); i += coresPerRow {
		end := i + coresPerRow
		if end > len(items) {
			end = len(items)
		}
		rows = append(rows, strings.Join(items[i:end], "  "))
	}
	return strings.Join(rows, "\n"+strings.Repeat(" ", 10))
}

func bar(percent float64, width int, cfg Config) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(percent / 100 * float64(width))
	if percent > 0 && filled == 0 {
		filled = 1
	}

	body := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	if cfg.NoColor {
		return body
	}

	fg := barColor(percent)
	return lipgloss.NewStyle().Foreground(fg).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(theme.SumiInk4).Render(strings.Repeat("░", width-filled))
}

func batteryBar(percent float64, width int, cfg Config) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(percent / 100 * float64(width))
	if percent > 0 && filled == 0 {
		filled = 1
	}

	body := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	if cfg.NoColor {
		return body
	}

	fg := batteryColor(percent)
	return lipgloss.NewStyle().Foreground(fg).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(theme.SumiInk4).Render(strings.Repeat("░", width-filled))
}

func barColor(percent float64) lipgloss.Color {
	switch {
	case percent > 90:
		return theme.WaveRed
	case percent > 70:
		return theme.CarpYellow
	default:
		return theme.SpringGreen
	}
}

func batteryColor(percent float64) lipgloss.Color {
	switch {
	case percent < 20:
		return theme.WaveRed
	case percent < 50:
		return theme.CarpYellow
	default:
		return theme.SpringGreen
	}
}

func trend(current, previous float64, cfg Config) string {
	const threshold = 2.0

	diff := current - previous
	if math.Abs(diff) < threshold {
		return style(cfg, theme.SumiInk4).Render("─")
	}
	if diff > 0 {
		return style(cfg, theme.SpringGreen).Render("▲")
	}
	return style(cfg, theme.WaveRed).Render("▼")
}

func renderProcesses(procs []stats.ProcessInfo, cfg Config) string {
	lines := []string{}
	for i, p := range procs {
		name := p.Name
		if len(name) > 18 {
			name = name[:18]
		}

		var line string
		if i == 0 {
			label := style(cfg, theme.OniViolet).Bold(true).Render(fmt.Sprintf("%-10s", "Top"))
			line = label + style(cfg, theme.FujiWhite).Render(fmt.Sprintf("%-20s%5.1f%%", name, p.CPUPercent))
		} else {
			line = strings.Repeat(" ", 10) + style(cfg, theme.FujiWhite).Render(fmt.Sprintf("%-20s%5.1f%%", name, p.CPUPercent))
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func palette(cfg Config) string {
	colors := []lipgloss.Color{
		theme.WaveRed,
		theme.SurimiOrange,
		theme.CarpYellow,
		theme.SpringGreen,
		theme.WaveAqua2,
		theme.SpringBlue,
		theme.CrystalBlue,
		theme.OniViolet,
	}

	blocks := make([]string, 0, len(colors))
	for _, c := range colors {
		if cfg.NoColor {
			blocks = append(blocks, "██")
			continue
		}
		blocks = append(blocks, lipgloss.NewStyle().Foreground(c).Render("██"))
	}
	return strings.Join(blocks, " ")
}

func bytes(value uint64) string {
	const unit = 1024
	if value < unit {
		return fmt.Sprintf("%d B", value)
	}

	div, exp := uint64(unit), 0
	for n := value / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(value)/float64(div), "KMGTPE"[exp])
}

func style(cfg Config, fg lipgloss.Color) lipgloss.Style {
	s := lipgloss.NewStyle()
	if !cfg.NoColor {
		s = s.Foreground(fg)
	}
	return s
}

func color(cfg Config, fg lipgloss.Color) lipgloss.Color {
	if cfg.NoColor {
		return ""
	}
	return fg
}

func RenderHelp(cfg Config) string {
	title := style(cfg, theme.CrystalBlue).Bold(true).Render("  keybindings")
	sep := style(cfg, theme.SumiInk4).Render("  ─────────────────────────")

	keys := []struct {
		key  string
		desc string
	}{
		{"a", "cycle ascii art (pulse > os > none)"},
		{"c", "toggle compact mode"},
		{"n", "toggle no-color mode"},
		{"p", "toggle top processes"},
		{"+/-", "adjust refresh interval"},
		{"?", "close this help"},
		{"q", "quit"},
	}

	lines := []string{"", title, sep}
	for _, k := range keys {
		keyStr := style(cfg, theme.SurimiOrange).Bold(true).Render(fmt.Sprintf("  %-6s", k.key))
		descStr := style(cfg, theme.FujiWhite).Render(k.desc)
		lines = append(lines, keyStr+descStr)
	}
	lines = append(lines, sep, "")

	return strings.Join(lines, "\n")
}
