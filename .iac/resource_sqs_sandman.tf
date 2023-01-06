
resource "aws_sqs_queue" "compse_sqs_sandman" {
  name                      = local.naming.prefix
  delay_seconds            = 90
  max_message_size          = 256000
  message_retention_seconds = 86400

  kms_master_key_id   = "alias/aws/sqs"
  kms_data_key_reuse_period_seconds = 300
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
