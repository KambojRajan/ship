package cmd

import (
	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new git repository",
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.Init(".")
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
