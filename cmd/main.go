package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	_sqs "serviceman/internal/adapters/framework/primary/sqs"
	"serviceman/internal/ports"
	"sync"
)

func main() {
	var SqsPollAdapter ports.SQSPORT
	sess := session.Must(session.NewSessionWithOptions(session.Options{Config: aws.Config{Region: aws.String("af-south-1")}}))
	sqsSvc := sqs.New(sess)
	queueUrl := "http://localhost:4566/000000000000/sandman-queue"
	SqsPollAdapter = _sqs.NewAdapter(sess, sqsSvc, queueUrl, ":8080")

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		SqsPollAdapter.RUN()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		chnMessages := make(chan *sqs.Message, 20)
		SqsPollAdapter.ReadMessages(chnMessages)
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
