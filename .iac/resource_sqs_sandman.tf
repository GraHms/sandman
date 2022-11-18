
resource "aws_sqs_queue" "compse_sqs_sandman" {
  name                      = local.naming.prefix
  delay_seconds            = 90
  max_message_size          = 2048
  message_retention_seconds = 86400
  receive_wait_time_seconds = 0
  visibility_timeout_seconds = 100
  redrive_policy = jsonencode({
    deadLetterTargetArn = aws_sqs_queue.compse_sqs_sandman_queue_deadletter.arn
    maxReceiveCount     = 100
  })

}

resource "aws_sqs_queue" "compse_sqs_sandman_queue_deadletter" {
  name = "${local.naming.prefix}-dlq"
}
