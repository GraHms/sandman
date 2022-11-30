package sqs

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"lab.dev.vm.co.mz/compse/sandman/internal/pkg/models"
)

type Adapter struct {
	SQS      sqsiface.SQSAPI
	QueueUrl string
	SqsSess  *session.Session
}

func NewAdapter(SqsSess *session.Session, sqs sqsiface.SQSAPI, queueUrl string) *Adapter {

	return &Adapter{
		SqsSess:  SqsSess,
		SQS:      sqs,
		QueueUrl: queueUrl,
	}
}
func (sqsa Adapter) SendMessage(body models.Body) error {
	strBody, err := sqsa.BodyComposer(body)
	if err != nil {
		return err
	}
	input := &sqs.SendMessageInput{
		DelaySeconds: aws.Int64(0),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Publisher": {
				DataType:    aws.String("String"),
				StringValue: aws.String("sandman"),
			},
			"Intent": {
				DataType:    aws.String("String"),
				StringValue: aws.String(body.Intent),
			},
			"Subscribers": {
				DataType:    aws.String("String"),
				StringValue: aws.String(body.Request.Sub),
			},
		},
		MessageBody: aws.String(strBody),
		QueueUrl:    aws.String(sqsa.QueueUrl),
	}
	message, err := sqsa.SQS.SendMessage(input)
	if err != nil {
		// could not send message to queue
		println("Could not send message to queue: ", err)
		return err
	}
	fmt.Printf("message sent to queue: %v", message.MessageId)
	return nil
}

func (sqsa Adapter) BodyComposer(body models.Body) (string, error) {
	data, err := json.Marshal(body)
	if err != nil {
		println("could not marshal sqs outbound body with traceId: ", body.TraceId)
		return "", err
	}
	return string(data), nil

}
