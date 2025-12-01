package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rollbackCmd)
	rollbackCmd.PersistentFlags().StringP("application", "a", "", "Application name to rollback")
	rollbackCmd.MarkPersistentFlagRequired("application")
	rollbackCmd.PersistentFlags().StringP("environment", "e", "", "Target environment for rollback")
	rollbackCmd.MarkPersistentFlagRequired("environment")
	rollbackCmd.PersistentFlags().IntP("task-definition", "t", 0, "Override the task definition version")
}

var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Redeploy the previous task definition",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Rollback functionality not yet implemented")
	},
}
