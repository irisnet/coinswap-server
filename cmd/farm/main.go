package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/irisnet/coinswap-server/config"
	"github.com/irisnet/coinswap-server/internal/app/server"
)

var (
	defaultCLIHome = os.ExpandEnv("$HOME/.farm")
	flagHome       = "home"
)

func main() {
	cobra.EnableCommandSorting = false
	rootCmd := &cobra.Command{
		Use:   "farm",
		Short: "farm Daemon (server)",
	}
	rootCmd.AddCommand(StartCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Failed executing mkr command: %s, exiting...\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

// StartCmd return the start command
func StartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Example: "farm start",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.Load(cmd, flagHome); err != nil {
				return err
			}
			server.Start()
			return nil
		},
	}
	cmd.Flags().String(flagHome, defaultCLIHome, "dapp server config path")
	return cmd
}
