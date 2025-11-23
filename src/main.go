package main

import (
	"os"

	"github.com/justmiles/ecs-deploy/src/cmd"
	lambda "github.com/justmiles/ecs-deploy/src/lambda"
)

func main() {

	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		lambda.Start()
	} else {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
		cmd.Execute()
	}
}
