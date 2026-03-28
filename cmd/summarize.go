package cmd

import (
	"os"

	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var summarizeCached bool

var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Suggest a commit message using AI based on the current diff",
	Long: `Summarize computes the current diff and uses the Claude API to suggest
a commit message in Conventional Commits format.

  ship summarize           # diff working tree vs index
  ship summarize --cached  # diff index vs HEAD (staged changes)

Requires ANTHROPIC_API_KEY to be set in the environment.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := os.Getwd()
		return commands.Summarize(path, summarizeCached)
	},
}

func init() {
	summarizeCmd.Flags().BoolVar(&summarizeCached, "cached", false, "Summarize staged changes (index vs HEAD)")
	rootCmd.AddCommand(summarizeCmd)
}
