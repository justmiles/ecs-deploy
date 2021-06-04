# ecs-deploy
A fast and flexible tool to deploy to Amazon Web Service's Elastic Container Service

    Usage:
      ecs-deploy [flags]
      ecs-deploy [command]

    Available Commands:
      help        Help about any command
      ship        Ship an application to ECS

    Flags:
      -d, --debug     Enable debug logging
      -h, --help      help for ecs-deploy
          --version   version for ecs-deploy

    Use "ecs-deploy [command] --help" for more information about a command.

## Installation

[Download the build for your machine](https://github.com/justmiles/ecs-deploy/releases), unzip, and add to your `$PATH`

## Usage
Run `ecs-deploy --help` to view available commands

Example:

    ecs-deploy ship --application myapp --environment qa --version latest

## Usage in AWS Lambda
Deployed this as a Lambda function and it can be invoked with the following JSON payload

  - Application - name of the application you want to update. 
  - Version - desired version
  - Environment - name of target environment (ECS Cluster)

You can use the included Terraform module to provision your Lambda function

1. Upload the [latest linux_amd64 release](https://github.com/justmiles/ecs-deploy/releases/latest) to S3

2. Implement the following module

        module "ecs_deploy" {
          source    = "git::ssh://git@github.com/justmiles/ecs-deploy.git?ref=master//tfmodule"
          s3_bucket = "my-s3-bucket"
          s3_key    = "ecs-deploy.zip"
        }
