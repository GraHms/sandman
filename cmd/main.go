package main

import (
	"crypto/tls"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"lab.dev.vm.co.mz/compse/sandman/config"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/application"
	_sqs "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/primary/sqs"
	_http "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/secundary/http"
	sqsOut "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/secundary/sqs"
	"lab.dev.vm.co.mz/compse/sandman/internal/ports"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	logger := config.Logger()
	var SqsPollAdapter ports.SQSPORT
	var AppAdapter ports.APPPort
	var requestAdater ports.RequestPORT
	var sqsOutbound ports.SecSQSPORT
	var httpClient ports.HTTPClient

	ssl := &tls.Config{
		InsecureSkipVerify: true,
	}
	httpClient = &http.Client{
		Timeout: 120 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: ssl,
		},
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			//Credentials: credentials.NewStaticCredentials("root", "root", ""),
			//Endpoint:    aws.String("http://localhost:4566"),
			Region: aws.String("af-south-1")}}))
	sqsSvc := sqs.New(sess)
	queueUrl := os.Getenv("SQS_QUEUE_URL")
	//queueUrl := "http://localhost:4566/000000000000/sandman-q"

	sqsOutbound = sqsOut.NewAdapter(sess, sqsSvc, queueUrl)
	requestAdater = _http.NewAdapter(httpClient, sqsOutbound, http.NewRequest)
	AppAdapter = application.NewAdapter(requestAdater)
	SqsPollAdapter = _sqs.NewAdapter(sess, sqsSvc, AppAdapter, queueUrl, ":8080")
	config.Init()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("server is running")
		SqsPollAdapter.RUN()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		chnMessages := make(chan *sqs.Message, 20)
		go SqsPollAdapter.ReadMessages(chnMessages)
		logger.Info("sqs is running")
		SqsPollAdapter.HandleMessages(chnMessages)
	}()
	wg.Wait()

}
