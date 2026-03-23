package cmd

import (
	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var CommitCmd = &cobra.Command{
	Use:   "commit [message]",
	Short: "Commit changes to the repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		message := args[0]
		err := commands.Commit(message, repoBasePath)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(CommitCmd)
}
