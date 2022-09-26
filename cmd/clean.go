package dudu

import (
	"github.com/spf13/cobra"
)

func init() {
	cleanCmd.Flags().BoolP("all", "a", false, "all runs")
	// duduCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean <path>",
	Short: "Show the difference between this run and the last",
	Run: func(cmd *cobra.Command, args []string) {
		// allFlag, _ := cmd.Flags().GetBool("all")
		// fmt.Printf("%v\n", allFlag)
	},
}
