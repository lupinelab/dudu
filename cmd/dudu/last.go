package dudu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dudu "git.lupinelab.co.uk/lupinelab/dudu/internal"
	"github.com/spf13/cobra"
)

var lastCmd = &cobra.Command{
	Use:   "last [path]",
	Short: "Show the difference between this run and the last",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get absolute filepath
		filePath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Check if there has been a run we can compare against
		compareTarget := "last"
		if _, err := os.Stat(dudu.TempDir + "/dudu" + strings.ReplaceAll(filePath, "/", ".") + ".last"); os.IsNotExist(err) {
			if _, err := os.Stat(dudu.TempDir + "/dudu" + strings.ReplaceAll(filePath, "/", ".") + ".first"); os.IsNotExist(err) {
				fmt.Println("No previous run found, please run \"dudu <path>\"")
				return
			} else {
				compareTarget = "first"
			}
		}

		// Run the dudu
		rawDuThisRun, err := dudu.Du(filePath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// convert rawDuThisRun to map[string]int
		duduThisRun := dudu.ParseDuData(rawDuThisRun)

		// convert compareTarget rawDu to map[string]int
		rawDuCompareTarget, err := os.ReadFile(dudu.TempDir + "/dudu" + strings.ReplaceAll(filePath, "/", ".") + "." + compareTarget)
		if err != nil {
			fmt.Println(err.Error())
		}
		duduCompareTarget := dudu.ParseDuData(rawDuCompareTarget)

		// Print comparison
		dudu.PrintDuduComparison(cmd, filePath, duduThisRun, duduCompareTarget)

		// Make file to write output into
		rawDuLastRun, err := os.Create(dudu.TempDir + "/dudu" + strings.ReplaceAll(filePath, "/", ".") + ".last")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rawDuLastRun.Close()

		// Write rawDuThisRun to file
		_, err = rawDuLastRun.WriteString(string(rawDuThisRun))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	},
}
