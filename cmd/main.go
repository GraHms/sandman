package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	_sqs "serviceman/internal/adapters/framework/primary/sqs"
	"serviceman/internal/ports"
	"sync"
)

func main() {
	var SqsPollAdapter ports.SQSPORT
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("af-south-1")}}))

	sqsSvc := sqs.New(sess)
	queueUrl := os.Getenv("SQS_QUEUE_URL")
	SqsPollAdapter = _sqs.NewAdapter(sess, sqsSvc, queueUrl, ":8080")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		println("server is running")
		SqsPollAdapter.RUN()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		chnMessages := make(chan *sqs.Message, 20)
		go SqsPollAdapter.ReadMessages(chnMessages)
		println("sqs is running")
		rwg := sync.WaitGroup{}
		for message := range chnMessages {
			rwg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				SqsPollAdapter.HandleMessage(message)
				SqsPollAdapter.DeleteMessage(message)
			}(&rwg)
			rwg.Done()
		}
	}()
	wg.Wait()

}
