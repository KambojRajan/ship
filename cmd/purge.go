package cmd

import (
	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var PurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Removes the unwanted temp files",
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.Purge(repoBasePath)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}
