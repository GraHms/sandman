package application

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-playground/assert/v2"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/mocks"
	outboundHttp "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/secundary/http"
	sqsOut "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/secundary/sqs"
	"lab.dev.vm.co.mz/compse/sandman/internal/ports"
	"net/http"
	"testing"
)

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

var errBody = `{
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
		"successStatus": 201,
		"hasCallBack": true,
		"callbacks": [
			{
				"description": "Payment Success feedback to commisison engine",
				"whenRequestStatus": 205,
				"successStatus": 205,
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
				"successStatus": 201,
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
var sess = session.Must(session.NewSessionWithOptions(session.Options{
	Config: aws.Config{}}))

func TestNewAdapter(t *testing.T) {
	var app ports.APPPort
	var sqsOutInbound ports.SecSQSPORT
	var reqAdpter ports.RequestPORT
	sqsOutInbound = sqsOut.NewAdapter(sess, mocks.GetMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	reqAdpter = outboundHttp.NewAdapter(client, sqsOutInbound, http.NewRequest)

	app = NewAdapter(reqAdpter)

	msg := sqs.Message{Body: &sqsBody}
	err := app.ProcessMessage(&msg)
	assert.Equal(t, err, nil)

}

func TestUnmarshalErr(t *testing.T) {
	var app ports.APPPort
	var sqsOutInbound ports.SecSQSPORT
	var reqAdpter ports.RequestPORT
	sqsOutInbound = sqsOut.NewAdapter(sess, mocks.GetMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	reqAdpter = outboundHttp.NewAdapter(client, sqsOutInbound, http.NewRequest)

	app = NewAdapter(reqAdpter)
	dummyBody := `foo:bar`
	msg := sqs.Message{Body: &dummyBody}
	err := app.ProcessMessage(&msg)
	assert.Equal(t, err, nil)

}

func TestRequestErr(t *testing.T) {
	var app ports.APPPort
	var sqsOutInbound ports.SecSQSPORT
	var reqAdpter ports.RequestPORT
	sqsOutInbound = sqsOut.NewAdapter(sess, mocks.GetMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	reqAdpter = outboundHttp.NewAdapter(client, sqsOutInbound, http.NewRequest)

	app = NewAdapter(reqAdpter)

	msg := sqs.Message{Body: &errBody}
	err := app.ProcessMessage(&msg)
	assert.NotEqual(t, err, nil)

}
