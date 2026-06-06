package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pranavray/ikiru/internal/ui"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	cfg := ui.Config{
		Interval:   time.Second,
		ASCIIStyle: ui.ASCIIStylePulse,
	}

	var intervalSeconds int

	rootCmd := &cobra.Command{
		Use:   "ikiru",
		Short: "A live neofetch-style system vitals TUI.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if intervalSeconds <= 0 {
				return fmt.Errorf("--interval must be greater than 0")
			}

			if cfg.ASCIIStyle != ui.ASCIIStylePulse && cfg.ASCIIStyle != ui.ASCIIStyleOS && cfg.ASCIIStyle != ui.ASCIIStyleNone {
				return fmt.Errorf("--ascii must be one of: pulse, os, none")
			}

			cfg.Interval = time.Duration(intervalSeconds) * time.Second

			if cfg.JSON {
				snap, err := ui.RenderOnceSnapshot(cfg)
				if err != nil {
					return err
				}
				data, err := json.MarshalIndent(snap, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(data))
				return nil
			}

			if cfg.Once {
				view, err := ui.RenderOnce(cfg)
				fmt.Println(view)
				return err
			}

			_, err := tea.NewProgram(ui.NewModel(cfg), tea.WithAltScreen()).Run()
			return err
		},
	}

	rootCmd.Flags().IntVar(&intervalSeconds, "interval", 1, "refresh interval in seconds")
	rootCmd.Flags().BoolVar(&cfg.Once, "once", false, "render a static snapshot and exit")
	rootCmd.Flags().BoolVar(&cfg.NoColor, "no-color", false, "strip all terminal colors")
	rootCmd.Flags().BoolVar(&cfg.Compact, "compact", false, "render condensed stats without ASCII art")
	rootCmd.Flags().StringVar(&cfg.ASCIIStyle, "ascii", ui.ASCIIStylePulse, "ASCII art style: pulse, os, none")
	rootCmd.Flags().BoolVar(&cfg.JSON, "json", false, "output snapshot as JSON and exit")

	return rootCmd
}
