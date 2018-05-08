package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/justmiles/ecs-deploy/src/deployer"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	application string
	version     string
	environment string
	noWait      bool
)

func init() {
	rootCmd.AddCommand(shipCmd)

	shipCmd.Flags().StringVarP(&application, "application", "a", "", "Application name to deploy")
	shipCmd.MarkFlagRequired("application")

	shipCmd.Flags().StringVarP(&version, "version", "v", "", "Desired version of application")
	shipCmd.MarkFlagRequired("version")

	shipCmd.Flags().StringVarP(&environment, "environment", "e", "", "Target environment for deployment")
	shipCmd.MarkFlagRequired("environment")

	shipCmd.Flags().BoolVarP(&noWait, "no-wait", "w", false, "Deploy and exit; Do not wait for service to reach stable state")

}

var shipCmd = &cobra.Command{
	Use:   "ship",
	Short: "Ship an application to ECS",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Printf("Deploying %s@%s to %s\n", application, version, environment)

		results, err := deployer.PerformDeployment(deployer.DeploymentOptions{
			Application: application,
			Version:     version,
			Environment: environment,
			Description: "Desired version set by ecs-deploy CLI",
		})

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
				var sess = session.Must(session.NewSession())

				fmt.Println("Waiting for service to reach stable state")

				svc := ecs.New(sess)

				err := svc.WaitUntilServicesStable(&ecs.DescribeServicesInput{
					Cluster: aws.String(environment),
					Services: []*string{
						aws.String(application),
					},
				})

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			fmt.Printf("%s@%s successfully updated in %s\n", application, version, environment)
		} else {
			fmt.Printf("Error pushing updates to %s\n", environment)
		}

	},
}
