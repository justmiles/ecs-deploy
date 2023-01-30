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

### Homebrew

```bash
brew install chrispruitt/tap/ecs-deploy
```

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

1.  Upload the [latest linux_amd64 release](https://github.com/justmiles/ecs-deploy/releases/latest) to S3

2.  Implement the following module

        module "ecs_deploy" {
          source    = "git::ssh://git@github.com/justmiles/ecs-deploy.git?ref=master//tfmodule"
          s3_bucket = "my-s3-bucket"
          s3_key    = "ecs-deploy.zip"
        }

## Refreshing Secrets Strategy

In reference to the, `ecs ship` command, there is an optional `--refresh-secrets` flag. This is used to pull a list of ssm parameters based on the `--secrets-prefix`. It will update all container definitions secrets to match the result.

Additionally, there may be a need to have multiple container defnitions with unique secrets. For this, just organize the ssm parameter store with the specific container name appended to the secrets-prefix.

Consider the following example:

A task definition has two containers. A container with the name `aaa` and a container with the name `bbb`. We need each continer to have have a unique value for a `LOG_LEVEL` env var.

The ssm parameter store is configured like so:

```text
/prd/myapp/GLOBAL_SECRET  = foo
/prd/myapp/aaa/LOG_LEVEL  = WARN
/prd/myapp/bbb/LOG_LEVEL  = DEBUG
```

After running the command:

```bash
ecs-deploy ship -a myapp -e prd -v 1.0.0 --refresh-secrets
```

The secrets in the containers will result in the following:

```json
{
  "containerDefinitions": [
    {
      "name": "aaa",
      "secrets": [
        {
          "name": "GLOBAL_SECRET",
          "valueFrom": "arn:aws:ssm:us-east-1:11111111111:parameter/prd/myapp/GLOBAL_SECRET"
        },
        {
          "name": "LOG_LEVEL",
          "valueFrom": "arn:aws:ssm:us-east-1:11111111111:parameter/prd/myapp/aaa/LOG_LEVEL"
        }
      ]
    },
    {
      "name": "bbb",
      "secrets": [
        {
          "name": "GLOBAL_SECRET",
          "valueFrom": "arn:aws:ssm:us-east-1:11111111111:parameter/prd/myapp/GLOBAL_SECRET"
        },
        {
          "name": "LOG_LEVEL",
          "valueFrom": "arn:aws:ssm:us-east-1:11111111111:parameter/prd/myapp/bbb/LOG_LEVEL"
        }
      ]
    }
  ]
}
```
