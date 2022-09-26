package dudu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().BoolP("all", "a", false, "all paths")
	duduCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean <path>",
	Short: "Remove data from previous runs",
	Run: func(cmd *cobra.Command, args []string) {
		allFlag, _ := cmd.Flags().GetBool("all")
		if allFlag == true {
			files, _ := os.ReadDir(TempDir)
			for _, file := range files {
				os.Remove(TempDir + "/" + file.Name())
			}
		} else {
			filePaths, err := filepath.Glob(TempDir + "/dudu" + (strings.ReplaceAll(args[0], "/", ".")) + ".*")
			if err != nil {
				fmt.Println(err)
			}
			for _, file := range filePaths {
				os.Remove(file)
			}
		}
	},
}
