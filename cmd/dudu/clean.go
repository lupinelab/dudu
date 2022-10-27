package dudu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dudu "git.lupinelab.co.uk/lupinelab/dudu/internal"
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().Bool("all", false, "all records")
}

var cleanCmd = &cobra.Command{
	Use:   "clean [path]",
	Short: "Remove records of previous runs",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if we have an path or flag
		allFlag, _ := cmd.Flags().GetBool("all")
		if len(args) < 1 && !allFlag {
			fmt.Println("Please specify a run path to clean or use \"dudu clean --all\" to clean all records")
			return
		}

		// Remove all run records
		count := 0
		if allFlag {
			files, _ := os.ReadDir(dudu.TempDir)
			for _, file := range files {
				err := os.Remove(dudu.TempDir + "/" + file.Name())
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				count++
			}
			fmt.Printf("Removed %v records\n", count)
		} else {
			// Remove all runs records for path
			runPath, err := filepath.Abs(args[0])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			filePaths, err := filepath.Glob(dudu.TempDir + "/dudu" + (strings.ReplaceAll(runPath, "/", ".")) + ".*")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, file := range filePaths {
				os.Remove(file)
				count++
			}
			// Print the amount of records removed
			fmt.Printf("Removed %v records\n", count)
		}
	},
}
