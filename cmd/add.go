package cmd

import (
	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "AddIndex file contents to the index",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := commands.Add(args...); err != nil {
			// Print any errors that occur during staging
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
