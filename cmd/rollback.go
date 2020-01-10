package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.PersistentFlags().StringP("application", "a", "YOUR NAME", "Author name for copyright attribution")
	rollbackCmd.PersistentFlags().IntP("task-definition", "t", 0, "Override the task definition version")
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Redeploy the previous task definition",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hugo Static Site Generator v0.9 -- HEAD")
	},
}
