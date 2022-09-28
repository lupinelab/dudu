package dudu

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var TempDir = os.TempDir() + "/dudu"

func Du(arg string) ([]byte, error) {
	// Run the du command
	if arg[len(arg)-1:] != "/" {
		arg = arg + "/"
	}
	duCmd := "du --summarize --exclude=proc --exclude=sys --exclude=dev --exclude=run --exclude=home --one-file-system " + arg + "*"
	du := exec.Command("/bin/sh", "-c", duCmd)
	rawDu, err := du.CombinedOutput()

	return rawDu, err
}

func ParseDuData(rawDuData []byte) map[string]int {
	duData := map[string]int{}
	duSlice := strings.Fields(string(rawDuData))

	for s := 1; s < len(duSlice); s += 2 {
		dirSize, _ := strconv.Atoi(duSlice[s-1])
		duData[duSlice[s]] = dirSize
	}

	return duData
}

func PrintDuduComparison(cmd *cobra.Command, filePath string, thisRun map[string]int, comparisonRun map[string]int) {
	// get suffix for comparison run file
	suffix := cmd.Name()
	if cmd.Name() == "total" {
		suffix = "first"
	}

	// get comparison file
	compareFilePath, err := filepath.Abs(TempDir + "/dudu" + strings.ReplaceAll(filePath, "/", ".", ) + "." + suffix)
	if err != nil {
		fmt.Println(err.Error())
	}
	compareFile, err := os.Stat(compareFilePath)

	// Print header
	fmt.Printf("%v - %v (since %v)\n", strings.ToUpper(cmd.Name()), filePath, compareFile.ModTime())

	// Print Comparison
	keys := make([]string, 0, len(thisRun))
	for k := range thisRun {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		size := strconv.Itoa(thisRun[k])
		change := strconv.Itoa((thisRun[k] - comparisonRun[k]))
		mFlag, _ := cmd.Flags().GetBool("mebibytes")
		if mFlag == true {
			size = strconv.Itoa(thisRun[k]/1024) + "M"
			change = strconv.Itoa((thisRun[k]-comparisonRun[k])/1024) + "M"
		}
		space := strings.Repeat(" ", 20-len(size))
		fmt.Printf("%v%v%v", size, space, k)
		if change[0:1] == "-" {
			fmt.Printf(" %v\n", color.GreenString(change))
		} else if change[0:len(change)-1] == "0" {
			fmt.Printf(" %v\n", color.YellowString(change))
		} else {
			fmt.Printf(" %v\n", color.RedString(change))
		}
	}
}
