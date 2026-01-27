package cmd

import (
	"os"

	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [path]",
	Short: "Show the working tree status",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		} else {
			var err error
			path, err = os.Getwd()
			if err != nil {
				cmd.PrintErr(err)
				return
			}
		}
		commands.Status(path)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
