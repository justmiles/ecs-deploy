package deployer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// DeploymentOptions set the desired state of your deployment
type DeploymentOptions struct {
	// Application name as it exists in ECS
	Application string `json:"Application"`
	// Desired version of ECS Application
	Version string `json:"Version"`
	// Environment (ECS Cluster Name) you would like to deploy to
	Environment string `json:"Environment"`
	// Description is an optional parameter adding context to the change
	Description string `json:"Description"`
	// Role is the IAM role to use when invoking a deployment.
	Role string `json:"Role"`
	// MaxAttempts is the Number of attempts to wait for service to become stable, with subsequent 15 sec pause.
	MaxAttempts int `json:"MaxAttempts"`
	// RefreshSecrets will update all container definition secrets to include ssm parameters that exists with the prefix "/<cluster>/service/*"
	RefreshSecrets bool `json:"RefreshSecrets"`
	// The ssm parameter store prefix to pull secrets from. Default: "/<cluster>/service/*"
	SecretsPrefix []string `json:"SecretsPrefix"`
	// DryRun will preview changes
	DryRun bool `json:"DryRun"`
}

// DeploymentResults maintain the depyments latest results
type DeploymentResults struct {
	// SuccessfullyInvoked bool value depicting a successful deployment invocation
	SuccessfullyInvoked bool   `json:"SuccessfullyInvoked"`
	ClusterArn          string `json:"ClusterArn"`
	ServiceArn          string `json:"ServiceArn"`
	ServiceName         string `json:"ServiceName"`
	TaskDefinition      string `json:"TaskDefinition"`
}

func (depOpts *DeploymentOptions) SetDeploymentOptionsByEcsServiceTags() error {
	// TODO: get ecs service tags that start with ecs-deploy:

	var ecsClient *ecs.ECS
	if depOpts.Role != "" {
		creds := stscreds.NewCredentials(sess, depOpts.Role)
		ecsClient = ecs.New(sess, &aws.Config{Credentials: creds})
	} else {
		ecsClient = ecs.New(sess)
	}

	ecsDescribeServicesOutput, err := ecsClient.DescribeServices(&ecs.DescribeServicesInput{
		Cluster:  aws.String(depOpts.Environment),
		Services: aws.StringSlice([]string{depOpts.Application}),
	})

	if err != nil {
		return fmt.Errorf("Unable to describe services: %v", err)
	}
	if len(ecsDescribeServicesOutput.Services) == 0 {
		return fmt.Errorf("ECS service not found.")
	}

	tagsOutput, err := ecsClient.ListTagsForResource(&ecs.ListTagsForResourceInput{
		ResourceArn: ecsDescribeServicesOutput.Services[0].ServiceArn,
	})
	if err != nil {
		return fmt.Errorf("Unable to get ecs service tags: %v", err)
	}

	r, _ := regexp.Compile("ecs-deploy:.*")
	for _, tag := range tagsOutput.Tags {
		if r.MatchString(*tag.Key) {
			switch strings.Split(*tag.Key, ":")[1] {
			case "refresh-secrets":
				value, err := strconv.ParseBool(*tag.Value)
				if err != nil {
					depOpts.RefreshSecrets = false
				} else {
					depOpts.RefreshSecrets = value
				}
				fmt.Println(fmt.Sprintf("ECS service tag found: \"%s=%s\". Setting --refresh-secrets to %t", *tag.Key, *tag.Value, depOpts.RefreshSecrets))

			case "secrets-prefix":
				value := strings.Split(*tag.Value, ":")
				depOpts.SecretsPrefix = value
				fmt.Println(fmt.Sprintf("ECS service tag found: \"%s=%s\". Setting --secrets-prefix to %v", *tag.Key, *tag.Value, depOpts.SecretsPrefix))
			}
		}
	}

	return nil
}
