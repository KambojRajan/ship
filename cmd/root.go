package cmd

import (
	"os"
	"slices"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/utils"
	"github.com/spf13/cobra"
)

var repoBasePath string

var rootCmd = &cobra.Command{
	Use:   "ship",
	Short: "ship is a mini git implementation",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path, _ := os.Getwd()
		var err error
		repoBasePath, err = utils.ShipHasBeenInitRecursive(path)
		if err != nil {
			return err
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if slices.Contains(utils.RunPurgeCommands, cmd.Name()) {
			commands.Purge(repoBasePath)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
