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
		oneline, _ := cmd.Flags().GetBool("oneline")
		err := commands.Pour(path, oneline)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	PourCmd.Flags().Bool("oneline", false, "Show each commit on a single line")
	rootCmd.AddCommand(PourCmd)
}
