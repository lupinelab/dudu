package dudu

import (
	"fmt"
	"os"

	dudu "git.lupinelab.co.uk/jedrw/dudu/pkg"
	"github.com/spf13/cobra"
)

var duduCmd = &cobra.Command{
	Use:   "dudu",
	Short: "dudu shows the difference in size of each folder at the root (/) of a filesytem between each run or since the first run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check and/or create tempdir
		tempDir := os.TempDir() + "/dudu"
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			os.Mkdir(tempDir, 0640)
		}

		// Run the dudu
		out, err := dudu.Du(args[0])
		if err != nil {
			fmt.Println(err)

		}

		// Make file to write output
		duduFirst, err := os.Create(tempDir + "/dudu.first")
		if err != nil {
			fmt.Println(err)
		}
		defer duduFirst.Close()

		// Write output to file
		_, err = duduFirst.WriteString(out)
		if err != nil {
			fmt.Println(err)
		}

		// Print output
		fmt.Print(out)
	},
}

// func init() {
// 	duduCmd.Flags().BoolP("first", "f", false, "Compare against first run")
// 	duduCmd.Flags().BoolP("last", "l", false, "Compare against latest run")
// }

func Execute() error {
	return duduCmd.Execute()
}
