resource "aws_iam_role" "sandman_iam_role" {
  name               = local.iam_role_name
  assume_role_policy = <<EOF
{
 "Version": "2012-10-17",
 "Statement": [
   {
     "Action": "sts:AssumeRole",
     "Principal": {
       "Service": "ecs-tasks.amazonaws.com"
     },
     "Effect": "Allow",
     "Sid": ""
   }
 ]
}
EOF
}

resource "aws_iam_policy" "sandman_policy" {

  name        = "commission-engine-ecs-task-policy"
  path        = "/"
  description = "Allow ECS task to read SQS Queue, Get Secret Value and Decrypt KMS"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action : [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:ChangeMessageVisibility"
        ],
        "Effect" : "Allow",
        "Resource" : aws_sqs_queue.compse_sqs_sandman.arn
      },
    ]
  })
}


resource "aws_iam_role_policy_attachment" "iam_role_policy_attachemtn" {
  role       = aws_iam_role.sandman_iam_role.name
  policy_arn = aws_iam_policy.sandman_policy.arn
}