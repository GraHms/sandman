package ports

import "github.com/aws/aws-sdk-go/service/sqs"

type APPPort interface {
	ProcessMessage(chan *sqs.Message)
	DeleteMessage()
}
