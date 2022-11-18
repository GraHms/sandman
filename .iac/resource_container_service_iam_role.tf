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

  name        = local.iam_role_name
  path        = "/"
  description = "Allow ECS task to read SQS Queue"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
        {
            "Effect": "Allow",
            "Action": [
                "sqs:SendMessage",
                "sqs:DeleteMessage",
                "sqs:ChangeMessageVisibility",
                "sqs:ReceiveMessage",
                "sqs:GetQueueAttributes",
                "sqs:TagQueue",
                "sqs:UntagQueue",
                "sqs:PurgeQueue"
            ],
            "Resource": aws_sqs_queue.compse_sqs_sandman.arn
        }
    ]
  })
}


resource "aws_iam_role_policy_attachment" "iam_role_policy_attachemtn" {
  role       = aws_iam_role.sandman_iam_role.name
  policy_arn = aws_iam_policy.sandman_policy.arn
}
