package cmd

import (
	"os"

	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var PourCmd = &cobra.Command{
	Use:   "pour",
	Short: "Display the commit history",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := os.Getwd()
		err := commands.Pour(path)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(PourCmd)
}

