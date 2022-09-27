package dudu

import (
	"fmt"
	"os"
	"strings"

	dudu "git.lupinelab.co.uk/lupinelab/dudu/internal"
	"github.com/spf13/cobra"
)

// func init() {
// 	duduCmd.AddCommand(totalCmd)
// }

var totalCmd = &cobra.Command{
	Use:   "total [path]",
	Short: "Show the difference between this run and the first",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if there has been a run already
		if _, err := os.Stat(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".first"); os.IsNotExist(err) {
			fmt.Println("No previous run found, please run \"dudu <path>\"")
			return
		}

		// Run the dudu
		rawDuThisRun, err := dudu.Du(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

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

		// convert rawDuThisRun to map[string]int
		duduThisRun := dudu.ParseDuData(rawDuThisRun)

		// convert firstRun rawDu to map[string]int
		rawDuCompareTarget, err := os.ReadFile(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".first")
		duduCompareTarget := dudu.ParseDuData(rawDuCompareTarget)

		// Print comparision
		dudu.PrintDuduComparison(cmd, duduThisRun, duduCompareTarget)
	},
}
