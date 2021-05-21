package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of pd",
	Long:  `All software has versions. This is geit's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("geit Static Site Generator v0.9 -- HEAD")
	},
}
