package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "gcmd",
	Short: "gcmd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gcmd is running")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
