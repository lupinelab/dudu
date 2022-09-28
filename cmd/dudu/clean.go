package dudu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean [path]",
	Short: "Remove records of previous runs",
	Run: func(cmd *cobra.Command, args []string) {
		count := 0
		allFlag, _ := cmd.Flags().GetBool("all")
		if allFlag == true {
			files, _ := os.ReadDir(TempDir)
			for _, file := range files {
				err := os.Remove(TempDir + "/" + file.Name())
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				count++
			}
			fmt.Printf("Removed %v records\n", count)
		} else {
			if len(args) < 1 {
				fmt.Println("Please specify a run path to clean or use \"dudu clean --all\" to clean all records")
				return
			} else {
				filePaths, err := filepath.Glob(TempDir + "/dudu" + (strings.ReplaceAll(args[0], "/", ".")) + ".*")
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				for _, file := range filePaths {
					os.Remove(file)
					count++
				}
				fmt.Printf("Removed %v records\n", count)
			}
		}
	},
}
