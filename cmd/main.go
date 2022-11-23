package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"net/http"
	"os"
	"serviceman/internal/adapters/application"
	_sqs "serviceman/internal/adapters/framework/primary/sqs"
	_http "serviceman/internal/adapters/framework/secundary/http"
	sqsOut "serviceman/internal/adapters/framework/secundary/sqs"
	"serviceman/internal/ports"
	"sync"
)

func main() {
	var SqsPollAdapter ports.SQSPORT
	var AppAdapter ports.APPPort
	var requestAdater ports.RequestPORT
	var sqsOutbound ports.SecSQSPORT
	httpClient := http.Client{}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			//Credentials: credentials.NewStaticCredentials("root", "root", ""),
			//Endpoint:    aws.String("http://localhost:4566"),
			Region: aws.String("af-south-1")}}))
	sqsSvc := sqs.New(sess)
	queueUrl := os.Getenv("SQS_QUEUE_URL")
	//queueUrl := "http://localhost:4566/000000000000/sandman-q"

	sqsOutbound = sqsOut.NewAdapter(sess, sqsSvc, queueUrl)
	requestAdater = _http.NewAdapter(&httpClient, sqsOutbound)
	AppAdapter = application.NewAdapter(requestAdater)
	SqsPollAdapter = _sqs.NewAdapter(sess, sqsSvc, AppAdapter, queueUrl, ":8080")

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
		SqsPollAdapter.HandleMessages(chnMessages)
	}()
	wg.Wait()

}
