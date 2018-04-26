# ecs-deploy
Lambda function to bump image version in ECS

## Usage
After you've deployed this as a Lambda function it can be invoked with the following JSON payload

  - Application - name of the application you want to update. 
  - Version - desired version
  - Environment - name of target environment (ECS Cluster)

Note: This function expects you're following the practice that an ECS Service is in the format `<env>-<application>`

## Terraform Module

1. Upload the [latest release](https://github.com/justmiles/ecs-deploy/releases/latest) to S3

2. Implement the following module 

        module "ecs_deploy" {
          source    = "git::ssh://git@github.com/justmiles/ecs-deploy.git?ref=master//tfmodule"
          s3_bucket = "my-s3-bucket"
          s3_key    = "LAMBDA_ARTIFACTS/ecs-deploy/ecs-deploy-v0.1.0.zip"
        }
