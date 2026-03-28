package cmd

import (
	"os"

	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent <description>",
	Short: "Stage all changes, generate a commit message with AI, and commit",
	Long: `Agent stages all working-tree changes, sends the diff plus your description
to the Claude API, generates an appropriate commit message, and commits automatically.

  ship agent "add retry logic to HTTP client"
  ship agent "fix null pointer in auth middleware"

Requires ANTHROPIC_API_KEY to be set in the environment.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := os.Getwd()
		return commands.Agent(args[0], path)
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
