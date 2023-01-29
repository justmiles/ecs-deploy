package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/justmiles/ecs-deploy/src/deployer"
	"github.com/spf13/cobra"
)

var (
	noWait            bool
	deploymentOptions = deployer.DeploymentOptions{
		Description: "Desired version set by ecs-deploy CLI",
	}
)

func init() {
	rootCmd.AddCommand(shipCmd)

	shipCmd.Flags().StringVarP(&deploymentOptions.Application, "application", "a", "", "Application name to deploy")
	shipCmd.MarkFlagRequired("application")

	shipCmd.Flags().StringVarP(&deploymentOptions.Version, "version", "v", "", "Desired version of application")
	shipCmd.MarkFlagRequired("version")

	shipCmd.Flags().StringVarP(&deploymentOptions.Environment, "environment", "e", "", "Target environment for deployment")
	shipCmd.MarkFlagRequired("environment")

	shipCmd.Flags().StringVarP(&deploymentOptions.Role, "role", "r", "", "An IAM role ARN to assume before invoking a deployment.")

	shipCmd.Flags().IntVar(&deploymentOptions.MaxAttempts, "max-attempts", 40, "Number of attempts (with subsequent 15 sec pause) to wait for service to become stable")

	shipCmd.Flags().BoolVarP(&noWait, "no-wait", "w", false, "Deploy and exit; Do not wait for service to reach stable state")

	shipCmd.Flags().BoolVar(&deploymentOptions.RefreshSecrets, "refresh-secrets", false, "Replace task defintion secrets with all ssm paramters with a prefix matching the 'secrets-prefix'")

	shipCmd.Flags().BoolVar(&deploymentOptions.DryRun, "dry-run", false, "Show changes without modifying resources.")

	shipCmd.Flags().StringVarP(&deploymentOptions.SecretsPrefix, "secrets-prefix", "p", "", "The ssm parameter store prefix to pull secrets from. Default: \"/<environment>/<application>/\"")
}

var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Ship an application to ECS",
	Run: func(cmd *cobra.Command, args []string) {

		if deploymentOptions.SecretsPrefix == "" {
			deploymentOptions.SecretsPrefix = fmt.Sprintf("/%s/%s", deploymentOptions.Environment, deploymentOptions.Application)
		}

		fmt.Printf("Deploying %s@%s to %s\n", deploymentOptions.Application, deploymentOptions.Version, deploymentOptions.Environment)
		results, err := deployer.PerformDeployment(deploymentOptions)
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

			fmt.Printf("%s@%s successfully updated in %s\n", deploymentOptions.Application, deploymentOptions.Version, deploymentOptions.Environment)
		} else {
			fmt.Printf("Error pushing updates to %s\n", deploymentOptions.Environment)
		}

	},
}
