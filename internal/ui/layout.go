package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pranavray/ikiru/internal/theme"
)

func RenderLayout(snapshot Snapshot, cfg Config) string {
	stats := renderStats(snapshot, cfg)
	if cfg.Compact || cfg.ASCIIStyle == ASCIIStyleNone {
		return stats
	}

	art := ASCIIArt(cfg.ASCIIStyle)
	artHeight := lipgloss.Height(art)
	statsHeight := lipgloss.Height(stats)
	if artHeight < statsHeight {
		art = strings.Repeat("\n", (statsHeight-artHeight)/2) + art
	}

	left := style(cfg, theme.SpringBlue).Render(art)
	right := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(color(cfg, theme.SumiInk4)).
		PaddingLeft(3).
		Render(stats)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, strings.Repeat(" ", 4), right)
}

func renderStats(snapshot Snapshot, cfg Config) string {
	lines := []string{
		row("OS", snapshot.System.OS+"  "+snapshot.System.Kernel, cfg),
		row("Host", snapshot.System.Hostname, cfg),
		row("Uptime", snapshot.System.Uptime, cfg),
		row("Shell", snapshot.System.Shell, cfg),
		row("CPU", snapshot.CPU.Model, cfg),
		row("Cores", renderCoreGrid(snapshot.CPU.Cores, cfg), cfg),
		row("Memory", fmt.Sprintf("%s / %s  %s", bytes(snapshot.Memory.Used), bytes(snapshot.Memory.Total), bar(snapshot.Memory.UsedPercent, 14, cfg)), cfg),
		row("Disk", fmt.Sprintf("%s / %s  %s", bytes(snapshot.Disk.Used), bytes(snapshot.Disk.Total), bar(snapshot.Disk.UsedPercent, 14, cfg)), cfg),
		row("Network", fmt.Sprintf("↑ %s/s  ↓ %s/s", bytes(uint64(snapshot.Network.BytesOutPerSec)), bytes(uint64(snapshot.Network.BytesInPerSec))), cfg),
		row("Processes", fmt.Sprintf("%d", snapshot.System.Processes), cfg),
	}

	if snapshot.Battery != nil && snapshot.Battery.Present {
		lines = append(lines, row("Battery", fmt.Sprintf("%.0f%% %s", snapshot.Battery.Percent, snapshot.Battery.State), cfg))
	}

	lines = append(lines, "", palette(cfg))
	return strings.Join(lines, "\n")
}

func row(label string, value string, cfg Config) string {
	labelText := fmt.Sprintf("%-10s", label)
	return style(cfg, theme.OniViolet).Bold(true).Render(labelText) + style(cfg, theme.FujiWhite).Render(value)
}

func renderCoreGrid(cores []float64, cfg Config) string {
	if len(cores) == 0 {
		return "unavailable"
	}

	items := make([]string, 0, len(cores))
	for i, usage := range cores {
		items = append(items, fmt.Sprintf("c%d %s %3.0f%%", i, bar(usage, 6, cfg), usage))
	}

	rows := []string{}
	for i := 0; i < len(items); i += 4 {
		end := i + 4
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

	return lipgloss.NewStyle().Foreground(theme.SpringGreen).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(theme.SumiInk4).Render(strings.Repeat("░", width-filled))
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
