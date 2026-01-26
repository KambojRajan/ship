package cmd

import (
	"fmt"

	"github.com/KambojRajan/ship/commands"
	"github.com/KambojRajan/ship/core/utils"
	"github.com/spf13/cobra"
)

var (
	catFileFlagP bool // pretty-print
	catFileFlagT bool // tree
	catFileFlagC bool // commit
	catFileFlagS bool // size
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file <hash>",
	Short: "Print raw content of object",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hash := args[0]
		flag := utils.CatFileFormatPretty

		if catFileFlagP {
			flag = utils.CatFileFormatPretty
		} else if catFileFlagT {
			flag = utils.CatFileFormatTree
		} else if catFileFlagC {
			flag = utils.CatFileFormatCommit
		} else if catFileFlagS {
			flag = utils.CatFileContentSize
		}

		result, err := commands.CatFile(hash, flag)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Print(result)
		if len(result) > 0 && result[len(result)-1] != '\n' {
			fmt.Println()
		}
	},
}

func init() {
	catFileCmd.Flags().BoolVarP(&catFileFlagP, "pretty", "p", false, "Pretty-print the object content")
	catFileCmd.Flags().BoolVarP(&catFileFlagT, "tree", "t", false, "Show tree object")
	catFileCmd.Flags().BoolVarP(&catFileFlagC, "commit", "c", false, "Show commit object")
	catFileCmd.Flags().BoolVarP(&catFileFlagS, "size", "s", false, "Show object size")
	rootCmd.AddCommand(catFileCmd)
}
