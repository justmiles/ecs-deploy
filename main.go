package main

import (
	"os"

	"github.com/justmiles/ecs-deploy/cmd"
	"github.com/justmiles/ecs-deploy/lambda"
)

func main() {

	if os.Getenv("_LAMBDA_SERVER_PORT") != "" {
		ld.Start()
	} else {
		os.Setenv("AWS_SDK_LOAD_CONFIG", "true")
		cmd.Execute()
	}
}
