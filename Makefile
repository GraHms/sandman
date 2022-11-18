aws --endpoint-url=http://localhost:4566 sqs create-queue --queue-name sandman-queue

#sudo snap install aws-cli --classic
awslocal sqs send-message --queue-url "http://localhost:4566/000000000000/sandman-queue" --message-body test
{
    "MD5OfMessageBody": "098f6bcd4621d373cade4e832627b4f6",
    "MessageId": "74861aab-05f8-0a75-ae20-74d109b7a76e"
}
