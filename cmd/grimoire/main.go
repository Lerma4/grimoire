// Command grimoire is the entrypoint for the Grimoire terminal app.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Lerma4/grimoire/internal/app"
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

// runTUI opens the terminal interface. The full TUI is wired in a later step;
// until then this reports progress so the binary is still usable.
func runTUI() error {
	fmt.Println("Grimoire — TUI not yet wired in this build.")
	fmt.Println("Try: grimoire version | grimoire doctor")
	return nil
}
