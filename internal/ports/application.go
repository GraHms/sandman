package ports

import "github.com/aws/aws-sdk-go/service/sqs"

type APPPort interface {
	ProcessMessage(*sqs.Message) error
	//DeleteMessage()

}
