package cmd

import (
	"os"

	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes to the repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]
		path, _ := os.Getwd()
		err := commands.Commit(message, path)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(CommitCmd)
}
