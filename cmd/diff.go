package cmd

import (
	"os"

	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var diffCached bool

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show changes between working tree and index, or index and HEAD",
	Long: `Show changes between the working tree and the index (default), or between
the index and the HEAD commit (--cached / --staged).

Examples:
  ship diff              # unstaged changes (working tree vs index)
  ship diff --cached     # staged changes (index vs HEAD)
  ship diff --staged     # alias for --cached`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := os.Getwd()
		if err := commands.Diff(path, diffCached); err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	diffCmd.Flags().BoolVar(&diffCached, "cached", false, "Compare index against HEAD instead of working tree")
	diffCmd.Flags().BoolVar(&diffCached, "staged", false, "Alias for --cached")
	rootCmd.AddCommand(diffCmd)
}
