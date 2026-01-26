package cmd

import (
	"fmt"

	"github.com/KambojRajan/ship/commands"
	"github.com/spf13/cobra"
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Print raw content of object",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := commands.CatFile(args)
		if err != nil {
			_ = fmt.Errorf(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
}
