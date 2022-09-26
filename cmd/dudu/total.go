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
	duduCmd.AddCommand(totalCmd)
}

var totalCmd = &cobra.Command{
	Use:   "total <path>",
	Short: "Show the difference between this run and the first",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check if there has been a run already
		if _, err := os.Stat(TempDir + "/dudu." + strings.ReplaceAll(args[0], "/", ".") + ".first"); os.IsNotExist(err) {
			fmt.Println("No previous run found, please run \"dudu <path>\"")
			return
		}

		// Run the dudu
		duRawThisRun, err := dudu.Du(args[0])
		if err != nil {
			fmt.Println(err)
		}

		// Make file to write output into
		duduLastFile, err := os.Create(TempDir + "/dudu." + strings.ReplaceAll(args[0], "/", ".") + ".last")
		if err != nil {
			fmt.Println(err)
		}
		defer duduLastFile.Close()

		// Write output to file
		_, err = duduLastFile.WriteString(string(duRawThisRun))

		// convert rawDuThisRun to map[string]int
		duDataThisRun := dudu.ParseDuData(duRawThisRun)

		// convert firstRun rawDu to map[string]int
		rawDuFirstRun, err := os.ReadFile(TempDir + "/dudu." + strings.ReplaceAll(args[0], "/", ".") + ".first")
		duDataFirstRun := dudu.ParseDuData(rawDuFirstRun)

		// Print comparision
		keys := make([]string, 0, len(duDataThisRun))
		for k := range duDataThisRun {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			space := strings.Repeat(" ", 20-len(strconv.Itoa(duDataThisRun[k])))
			change := strconv.Itoa((duDataThisRun[k]-duDataFirstRun[k])/1024) + "M"
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
