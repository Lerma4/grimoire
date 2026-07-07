// Command grimoire is the entrypoint for the Grimoire terminal app.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Lerma4/grimoire/internal/app"
	"github.com/Lerma4/grimoire/internal/store"
	"github.com/Lerma4/grimoire/internal/tui"
)

func main() {
	if err := newRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "grimoire:", err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "grimoire",
		Short: "A terminal grimoire for tasks and markdown notes",
		Long: "Grimoire is a local-first terminal app for managing tasks and\n" +
			"markdown notes, with Vim-style keybindings and task↔note linking.",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
	root.AddCommand(newTUICmd(), newVersionCmd(), newDoctorCmd())
	return root
}

func newTUICmd() *cobra.Command {
	return &cobra.Command{
		Use:           "tui",
		Short:         "Open the terminal user interface",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTUI()
		},
	}
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the Grimoire version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), "grimoire", app.Version)
		},
	}
}

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "doctor",
		Short:         "Check database, config paths and terminal environment",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := app.DefaultConfig()
			if err != nil {
				return err
			}
			for _, r := range app.Doctor(cfg) {
				mark := "ok"
				if !r.OK {
					mark = "FAIL"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  [%s] %-18s %s\n", mark, r.Name, r.Msg)
			}
			return nil
		},
	}
}

// runTUI opens and migrates the database, then starts the terminal interface.
func runTUI() error {
	cfg, err := app.DefaultConfig()
	if err != nil {
		return err
	}
	db, err := store.Open(cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := store.Migrate(db); err != nil {
		return err
	}
	if err := store.SeedIfEmpty(db); err != nil {
		return err
	}

	p := tea.NewProgram(tui.NewModel(tui.Deps{DB: db, DBPath: cfg.DBPath}), tea.WithAltScreen())
	_, err = p.Run()
	return err
}
