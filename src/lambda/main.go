package ld

import (
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/justmiles/ecs-deploy/src/deployer"
)

var (
	// ErrInvalidInputProvided is thrown when the input does not include Application, Version, and Environment
	ErrInvalidInputProvided = errors.New("invalid input - ensure Application, Version, and Environment is provided")
	defaultSSMDescription   = "desired version set by lambda ecs-deploy"
)

// Handler is your Lambda function handler
// It uses the DeploymentOptions JSON event to invoke a deployment against ECS
func Handler(depOpts deployer.DeploymentOptions) (s string, err error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing request to deploy %s@%s to %s\n", depOpts.Application, depOpts.Version, depOpts.Environment)

	if !inputValidation(depOpts) {
		return s, ErrInvalidInputProvided
	}

	if depOpts.Description == "" {
		depOpts.Description = defaultSSMDescription
	}

	return deployer.PerformDeployment(depOpts)
}

// Start the lambda function
func Start() {
	lambda.Start(Handler)
}

func inputValidation(depOpts deployer.DeploymentOptions) bool {
	if depOpts.Application == "" {
		return false
	}

	if depOpts.Version == "" {
		return false
	}

	if depOpts.Environment == "" {
		return false
	}
	return true
}
