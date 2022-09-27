package dudu

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	dudu "git.lupinelab.co.uk/lupinelab/dudu/internal"
	"github.com/spf13/cobra"
)

func init() {
	cobra.EnableCommandSorting = false
	duduCmd.AddCommand(lastCmd)
	duduCmd.AddCommand(totalCmd)
	duduCmd.AddCommand(cleanCmd)
	duduCmd.PersistentFlags().BoolP("mebibytes", "m", false, "Print sizes in MiB")
	duduCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	duduCmd.PersistentFlags().Lookup("help").Hidden = true
	duduCmd.CompletionOptions.DisableDefaultCmd = true
}

var TempDir = os.TempDir() + "/dudu"

var duduCmd = &cobra.Command{
	Use:   "dudu [path]",
	Short: "dudu shows the difference in size of each folder at the specified path between each run or since the first run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Run the dudu
		rawDu, err := dudu.Du(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Check/create tempdir
		if _, err := os.Stat(TempDir); os.IsNotExist(err) { //use errors.Is?
			os.Mkdir(TempDir, 0640)
		}

		// Make file to write output into
		filePath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		duduFirstFile, err := os.Create(TempDir + "/dudu" + strings.ReplaceAll(filePath, "/", ".") + ".first")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer duduFirstFile.Close()

		// Write output to file
		_, err = duduFirstFile.WriteString(string(rawDu))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// convert rawDu to map[string]int
		thisduData := dudu.ParseDuData(rawDu)

		// Print output
		keys := make([]string, 0, len(thisduData))
		for k := range thisduData {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		fmt.Println("FIRST\n=====")
		for _, k := range keys {
			size := strconv.Itoa(thisduData[k])
			mFlag, _ := cmd.Flags().GetBool("mebibytes")
			if mFlag == true {
				size = strconv.Itoa(thisduData[k]/1024) + "M"
			}
			space := strings.Repeat(" ", 20-len(size))
			fmt.Printf("%v%v%v\n", size, space, k)
		}
	},
}

func Execute() error {
	return duduCmd.Execute()
}
