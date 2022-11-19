package ports

import "github.com/aws/aws-sdk-go/service/sqs"

type SQSPORT interface {
	PollMessages(chn chan<- *sqs.Message)
	HandleMessage(msg *sqs.Message) error
	DeleteMessage(msg *sqs.Message)
	ReadMessages(chnMessages chan *sqs.Message)
	HandleMessages(chnMessages <-chan *sqs.Message)
	RUN()
}
