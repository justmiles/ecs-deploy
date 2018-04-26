resource "aws_iam_role" "ecs_deploy" {
  name = "ecs-deploy"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_policy" "ecs_deploy" {
  name        = "ecs-deploy"
  path        = "/"
  description = "Execution basics for ECS Deploy"
  policy      = "${var.iam_policy}"
}

resource "aws_iam_role_policy_attachment" "ecs_deploy" {
  role       = "${aws_iam_role.ecs_deploy.name}"
  policy_arn = "${aws_iam_policy.ecs_deploy.arn}"
}

variable "iam_policy" {
  default = <<EOF
{
   "Version":"2012-10-17",
   "Statement":[
      {
         "Sid":"",
         "Effect":"Allow",
         "Action":[
            "logs:PutLogEvents",
            "logs:DescribeLogStreams",
            "logs:DescribeLogGroups",
            "logs:CreateLogStream",
            "logs:CreateLogGroup",
            "ecs:*",
            "ecr:List*",
            "ssm:Get*",
            "ssm:Put*",
            "ssm:List*",
            "iam:PassRole"
         ],
         "Resource":"*"
      }
   ]
}
EOF
}
