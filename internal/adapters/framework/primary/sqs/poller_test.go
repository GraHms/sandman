package sqs

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-playground/assert/v2"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/application"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/mocks"
	outboundHttp "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/secundary/http"
	sqsOut "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/secundary/sqs"
	"lab.dev.vm.co.mz/compse/sandman/internal/ports"
	"net/http"
	"testing"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	Config: aws.Config{}}))
var appAdapter ports.APPPort
var requestAdpter ports.RequestPORT
var sqsOutInbound ports.SecSQSPORT

func TestAdapter_DeleteMessage(t *testing.T) {
	sqsOutInbound = sqsOut.NewAdapter(sess, mocks.GetMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	requestAdpter = outboundHttp.NewAdapter(client, sqsOutInbound, http.NewRequest)
	appAdapter = application.NewAdapter(requestAdpter)
	adp := NewAdapter(sess, mocks.GetMockSQSClient(), appAdapter, "", "8080")
	adp.DeleteMessage(&sqs.Message{})

}

func TestAdapter_HandleMessage(t *testing.T) {
	sqsOutInbound = sqsOut.NewAdapter(sess, mocks.GetMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	requestAdpter = outboundHttp.NewAdapter(client, sqsOutInbound, http.NewRequest)
	appAdapter = application.NewAdapter(requestAdpter)
	adp := NewAdapter(sess, mocks.GetMockSQSClient(), appAdapter, "", "8080")
	body := ``
	sqsMsg := &sqs.Message{
		Body: &body,
	}
	err := adp.HandleMessage(sqsMsg)
	assert.Equal(t, err, nil)

}

func TestAdapter_PollMessages(t *testing.T) {
	sqsOutInbound = sqsOut.NewAdapter(sess, mocks.GetMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	requestAdpter = outboundHttp.NewAdapter(client, sqsOutInbound, http.NewRequest)
	appAdapter = application.NewAdapter(requestAdpter)
	adp := NewAdapter(sess, mocks.GetMockSQSClient(), appAdapter, "https://queue.amazonaws.com/80398EXAMPLE/MyQueue", "8080")
	body := `{"foo":"bar"}`
	qurl := "https://queue.amazonaws.com/80398EXAMPLE/MyQueue"
	_, _ = adp.SQS.SendMessage(&sqs.SendMessageInput{MessageBody: &body, QueueUrl: &qurl})
	chnMessages := make(chan *sqs.Message)
	go adp.ReadMessages(chnMessages)
	resp := <-chnMessages
	assert.Equal(t, resp.Body, body)

}

var sqsBody = `{
	"name": "commission-payment",
	"traceId": "xr-1234567890",
	"groupReference": "",
	"origin": "commission-engine",
	"sandmanVersion": "1.0.0",
	"intent": "Make an mpesa payment and give feedback",
	"description": "",
	"journey": "Commission Payment",
	"owner": "COMPSE Squad",
	"request": {
		"retries": 3,
		"method": "POST",
		"contentType": "application/json",
		"body": {
			"paymentId": "paymentId",
			"beneficiaryAccount": "848255237",
			"paymentAmount": 32
		},
		"sub": "https://commission-payments-backend.svc.dev.ind.vm.co.mz/api/v1/payments",
		"headers": ""
	},
	"response": {
		"successStatus": 200,
		"hasCallBack": true,
		"callbacks": [
			{
				"description": "Payment Success feedback to commisison engine",
				"whenRequestStatus": 200,
				"successStatus": 200,
				"retries": 3,
				"method": "POST",
				"contentType": "application/json",
				"body": {
					"status": "completed",
					"statusChangeReason": "completed"
				},
				"mapToBody": [
					{
						"queryField": "paymentRef",
						"targetField": "paymentReference"
					}
				],
				"sub": "https://commission-backend.svc.dev.ind.vm.co.mz/api/v1/payments/3/feedback",
				"headers": ""
			},
			{
				"description": "Payment Success feedback to commisison engine",
				"whenRequestStatus": 500,
				"successStatus": 200,
				"retries": 3,
				"method": "POST",
				"contentType": "application/json",
				"body": {
					"status": "failed",
					"statusChangeReason": "failed"
				},
				"mapToBody": [
					{
						"queryField": "reason",
						"targetField": "statusChangeReason"
					}
				],
				"sub": "https://commission-backend.svc.dev.ind.vm.co.mz/api/v1/payments/3/feedback",
				"headers": ""
			},
			{
				"description": "Payment Success feedback to commisison engine",
				"whenRequestStatus": 504,
				"successStatus": 200,
				"retries": 3,
				"method": "POST",
				"contentType": "application/json",
				"body": {
					"status": "failed",
					"statusChangeReason": "Gateway timedout"
				},

				"sub": "https://commission-backend.svc.dev.ind.vm.co.mz/api/v1/payments/3/feedback",
				"headers": ""
			}
		]
	}
}`

func TestAdapter_HandleMessages(t *testing.T) {
	sqsOutInbound = sqsOut.NewAdapter(sess, mocks.GetMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	requestAdpter = outboundHttp.NewAdapter(client, sqsOutInbound, http.NewRequest)
	appAdapter = application.NewAdapter(requestAdpter)
	adp := NewAdapter(sess, mocks.GetMockSQSClient(), appAdapter, "https://queue.amazonaws.com/80398EXAMPLE/MyQueue", "8080")

	qurl := "https://queue.amazonaws.com/80398EXAMPLE/MyQueue"
	_, err := adp.SQS.SendMessage(&sqs.SendMessageInput{MessageBody: &sqsBody, QueueUrl: &qurl})
	chnMessages := make(chan *sqs.Message)
	adp.ReadMessages(chnMessages)

	go adp.HandleMessages(chnMessages)
	msg := <-chnMessages
	fmt.Printf("%v", msg)
	close(chnMessages)
	assert.Equal(t, err, nil)

}
