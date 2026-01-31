package cmd

import (
	"os"

	"github.com/KambojRajan/ship/core/entities"
	"github.com/spf13/cobra"
)

var PourCmd = &cobra.Command{
	Use:   "pour",
	Short: "Pour this gives you the commit history.",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := os.Getwd()
		_, err := entities.LoadCommits(path)
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(PourCmd)
}
