package dudu

import (
	"fmt"
	"os"
	"strings"

	dudu "git.lupinelab.co.uk/lupinelab/dudu/internal"
	"github.com/spf13/cobra"
)

// func init() {
// 	duduCmd.AddCommand(lastCmd)
// }

var lastCmd = &cobra.Command{
	Use:   "last [path]",
	Short: "Show the difference between this run and the last",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		compareTarget := "last"
		// Check if there has been a run we can compare against
		if _, err := os.Stat(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".last"); os.IsNotExist(err) {
			if _, err := os.Stat(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".first"); os.IsNotExist(err) {
				fmt.Println("No previous run found, please run \"dudu <path>\"")
				return
			} else {
				compareTarget = "first"
			}
		}

		// Run the dudu
		rawDuThisRun, err := dudu.Du(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// convert rawDuThisRun to map[string]int
		duduThisRun := dudu.ParseDuData(rawDuThisRun)

		// convert compareTarget rawDu to map[string]int
		rawDuCompareTarget, err := os.ReadFile(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + "." + compareTarget)
		duduCompareTarget := dudu.ParseDuData(rawDuCompareTarget)

		// Make file to write output into
		rawDuLastRun, err := os.Create(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".last")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer rawDuLastRun.Close()

		// Write output to file
		_, err = rawDuLastRun.WriteString(string(rawDuThisRun))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Print comparison
		dudu.PrintDuduComparison(cmd, duduThisRun, duduCompareTarget)
	},
}
