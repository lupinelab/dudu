package dudu

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	dudu "git.lupinelab.co.uk/lupinelab/dudu/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	duduCmd.AddCommand(lastCmd)
}

var lastCmd = &cobra.Command{
	Use:   "last <path>",
	Short: "Show the difference between this run and the last",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		compareTarget := "last"
		// Check if there has been a run we can compare against
		if _, err := os.Stat(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".last"); os.IsNotExist(err) {
			if _, err := os.Stat(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".first"); os.IsNotExist(err) {
				fmt.Println("No previous run found, please run \"dudu <path>\"")
				return
			}
			compareTarget = "first"
		}

		// Run the dudu
		duRawThisRun, err := dudu.Du(args[0])
		if err != nil {
			fmt.Println(err)
		}

		// convert rawDuThisRun to map[string]int
		duDataThisRun := dudu.ParseDuData(duRawThisRun)

		// convert compareTarget rawDu to map[string]int
		rawDuCompareTarget, err := os.ReadFile(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + "." + compareTarget)
		duDataCompareTarget := dudu.ParseDuData(rawDuCompareTarget)

		// Make file to write output into
		duduLastFile, err := os.Create(TempDir + "/dudu" + strings.ReplaceAll(args[0], "/", ".") + ".last")
		if err != nil {
			fmt.Println(err)
		}
		defer duduLastFile.Close()

		// Write output to file
		_, err = duduLastFile.WriteString(string(duRawThisRun))
		if err != nil {
			fmt.Println(err)
		}

		// Print comparison
		keys := make([]string, 0, len(duDataThisRun))
		for k := range duDataThisRun {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		fmt.Println("LAST\n====")
		for _, k := range keys {
			space := strings.Repeat(" ", 20-len(strconv.Itoa(duDataThisRun[k])))
			change := strconv.Itoa((duDataThisRun[k]-duDataCompareTarget[k])/1024) + "M"
			fmt.Printf("%d%v%v", duDataThisRun[k], space, k)
			if change[0:1] == "-" {
				fmt.Printf(" %v\n", color.GreenString(change))
			} else if change[0:len(change)-1] == "0" {
				fmt.Printf(" %v\n", color.YellowString(change))
			} else {
				fmt.Printf(" %v\n", color.RedString(change))
			}
		}
	},
}
