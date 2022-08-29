package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/justmiles/ecs-deploy/src/deployer"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(restartCmd)

	restartCmd.Flags().StringVarP(&deploymentOptions.Application, "application", "a", "", "Application name to redeploy")
	restartCmd.MarkFlagRequired("application")

	restartCmd.Flags().StringVarP(&deploymentOptions.Environment, "environment", "e", "", "Target environment")
	restartCmd.MarkFlagRequired("environment")

	restartCmd.Flags().StringVarP(&deploymentOptions.Role, "role", "r", "", "An IAM role ARN to assume before invoking a deployment.")

	restartCmd.Flags().IntVar(&deploymentOptions.MaxAttempts, "max-attempts", 40, "Number of attempts (with subsequent 15 sec pause) to wait for service to become stable")

	restartCmd.Flags().BoolVarP(&noWait, "no-wait", "w", false, "Redeploy and exit; Do not wait for service to reach stable state")

}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "gracefully restart/redeploy an application",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Redeploying %s in %s\n", deploymentOptions.Application, deploymentOptions.Environment)
		results, err := deployer.PerformReDeployment(deploymentOptions)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if debugEnabled {
			fmt.Println(results)
		}

		var depRes deployer.DeploymentResults
		err = json.Unmarshal([]byte(results), &depRes)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if depRes.SuccessfullyInvoked {

			if !noWait {
				fmt.Println("Waiting for service to reach stable state")

				err := deployer.WaitForDeployment(deploymentOptions)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			fmt.Printf("%s successfully restarted in %s\n", deploymentOptions.Application, deploymentOptions.Environment)
		} else {
			fmt.Printf("Error pushing updates to %s\n", deploymentOptions.Environment)
		}

	},
}
