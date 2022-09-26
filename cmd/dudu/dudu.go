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

var TempDir = os.TempDir() + "/dudu"

var duduCmd = &cobra.Command{
	Use:   "dudu <path>",
	Short: "dudu shows the difference in size of each folder at the specified path between each run or since the first run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Run the dudu
		rawDu, err := dudu.Du(args[0])
		if err != nil {
			fmt.Println(err)
		}

		// Check/create tempdir
		if _, err := os.Stat(TempDir); os.IsNotExist(err) {
			os.Mkdir(TempDir, 0640)
		}

		// Make file to write output into
		filePath, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println(err)
		}
		duduFirstFile, err := os.Create(TempDir + "/dudu" + strings.ReplaceAll(filePath, "/", ".") + ".first")
		if err != nil {
			fmt.Println(err)
		}
		defer duduFirstFile.Close()

		// Write output to file
		_, err = duduFirstFile.WriteString(string(rawDu))
		if err != nil {
			fmt.Println(err)
		}

		// convert rawDu to map[string]int
		thisduData := dudu.ParseDuData(rawDu)

		// Print output
		keys := make([]string, 0, len(thisduData))
		for k := range thisduData {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		for _, k := range keys {
			space := strings.Repeat(" ", 20-len(strconv.Itoa(thisduData[k])))
			fmt.Printf("%d%v%v\n", thisduData[k], space, k)
		}
	},
}

func init() {
	duduCmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	duduCmd.PersistentFlags().Lookup("help").Hidden = true
	duduCmd.CompletionOptions.DisableDefaultCmd = true
}

func Execute() error {
	return duduCmd.Execute()
}
