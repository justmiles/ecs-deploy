resource "aws_lambda_function" "ecs_deploy" {
  function_name    = "ecs-deploy"
  s3_bucket        = "${var.s3_bucket}"
  s3_key           = "${var.s3_key}"
  role             = "${aws_iam_role.ecs_deploy.arn}"
  handler          = "ecs-deploy"
  source_code_hash = "${var.source_code_hash}"
  runtime          = "go1.x"
}

variable "s3_bucket" {
  type = "string"
}

variable "s3_key" {
  type = "string"
}

variable "source_code_hash" {
  type    = "string"
  default = ""
}
