package sqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"serviceman/internal/ports"
	"sync"
)

type Adapter struct {
	SQS        *sqs.SQS
	AppAdapter ports.APPPort
	router     *gin.Engine
	QueueUrl   string
	addr       string
	SqsSess    *session.Session
}

func NewAdapter(SqsSess *session.Session, sqs *sqs.SQS, appAdapter ports.APPPort, queueUrl string, addr string) *Adapter {
	router := gin.Default()
	return &Adapter{
		SqsSess:    SqsSess,
		AppAdapter: appAdapter,
		router:     router,
		SQS:        sqs,
		QueueUrl:   queueUrl,
		addr:       addr,
	}
}
func (sqsa *Adapter) PollMessages(chn chan<- *sqs.Message) {

	for {
		output, err := sqsa.SQS.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(sqsa.QueueUrl),
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(3),
		})

		if err != nil {
			fmt.Printf("failed to fetch sqs message %v ", err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}

func (sqsa *Adapter) HandleMessage(msg *sqs.Message) error {
	//intent := *msg.MessageAttributes["Intent"]
	//pub := *msg.MessageAttributes["Publisher"]
	//sub := *msg.MessageAttributes["Subscribers"]
	fmt.Println("Receiving a message")
	err := sqsa.AppAdapter.ProcessMessage(msg)
	if err != nil {
		println(err.Error())
		return err
	}
	return nil
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

func (sqsa *Adapter) ReadMessages(chnMessages chan *sqs.Message) {

	go sqsa.PollMessages(chnMessages)

}

func (sqsa *Adapter) HandleMessages(chnMessages <-chan *sqs.Message) {
	wg := sync.WaitGroup{}
	for message := range chnMessages {
		wg.Add(1)
		func() {
			defer wg.Done()
			// HandleMessage message is error, do not delete from queue
			err := sqsa.HandleMessage(message)
			if err != nil {
				return
			}
			// if processed, delete from queue
			sqsa.DeleteMessage(message)
		}()
	}
	wg.Wait()
}
