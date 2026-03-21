package cmd

import (
	"os"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/utils"
	"github.com/spf13/cobra"
)

var repoBasePath string

var PurgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Removes the unwanted temp files",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path, _ := os.Getwd()

		var err error
		repoBasePath, err = utils.ShipHasBeenInitRecursive(path)
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := commands.Purge(repoBasePath)
		if err != nil {
			cmd.PrintErr(err)
		}
	},
}
