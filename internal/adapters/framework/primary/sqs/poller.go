package sqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
)

type Adapter struct {
	SQS      *sqs.SQS
	router   *gin.Engine
	QueueUrl string
	addr     string
	SqsSess  *session.Session
}

func NewAdapter(SqsSess *session.Session, sqs *sqs.SQS, queueUrl string, addr string) *Adapter {
	router := gin.Default()
	return &Adapter{
		SqsSess:  SqsSess,
		router:   router,
		SQS:      sqs,
		QueueUrl: queueUrl,
		addr:     addr,
	}
}
func (sqsa *Adapter) PollMessages(chn chan<- *sqs.Message) {

	for {
		output, err := sqsa.SQS.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sqsa.QueueUrl),
			MaxNumberOfMessages: aws.Int64(2),
			WaitTimeSeconds:     aws.Int64(10),
		})

		if err != nil {
			fmt.Printf("failed to fetch sqs message %v ", err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}

func (sqsa *Adapter) HandleMessage(msg *sqs.Message) {
	fmt.Println("RECEIVING MESSAGE >>> ")
	fmt.Println(*msg.Body)
}

func (sqsa *Adapter) DeleteMessage(msg *sqs.Message) {
	message, err := sqsa.SQS.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sqsa.QueueUrl),
		ReceiptHandle: msg.ReceiptHandle,
	})

	if err != nil {
		return
	}
	println("sqs message deleted: ", message.String())
}

func (sqsa *Adapter) ReadMessages(chnMessages chan<- *sqs.Message) {

	//chnMessages := make(chan *sqs.Message, 20)

	go sqsa.PollMessages(chnMessages)

}
