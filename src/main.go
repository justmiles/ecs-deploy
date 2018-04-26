package main

import (
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	// ErrInvalidInputProvided is thrown when the input does not include Application, Version, and Environment
	ErrInvalidInputProvided = errors.New("invalid input - ensure Application, Version, and Environment is provided")
	sess                    = session.Must(session.NewSession())
	defaultDescription      = "desired version set by lambda ecs-deploy"
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

// Handler is your Lambda function handler
// It uses the DeploymentOptions JSON event to invoke a deployment against ECS
func Handler(depOpts DeploymentOptions) (s string, err error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing request to deploy %s@%s to %s\n", depOpts.Application, depOpts.Version, depOpts.Environment)

	if !inputValidation(depOpts) {
		return s, ErrInvalidInputProvided
	}

	if depOpts.Description == "" {
		depOpts.Description = defaultDescription
	}

	// Set the desired application version
	err = setDesiredVersion(depOpts)
	if err != nil {
		return s, err
	}

	return performDeployment(depOpts)
}

func main() {
	lambda.Start(Handler)
}
