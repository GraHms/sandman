package http

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/go-playground/assert/v2"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/mocks"
	sqsOut "lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/secundary/sqs"
	"lab.dev.vm.co.mz/compse/sandman/internal/pkg/models"
	"lab.dev.vm.co.mz/compse/sandman/internal/ports"
	"net/http"
	"testing"
)

func TestAdapter_ConvertBodyResponse(t *testing.T) {
	adpter := Adapter{}
	body := `{"fake":"body"}`
	fakeResp := mocks.MockBody(body)
	res := adpter.ConvertBodyResponse(fakeResp)
	assert.Equal(t, res, map[string]interface{}{"fake": "body"})
}

func TestAdapter_BodyMapper(t *testing.T) {
	adpter := Adapter{}
	body := map[string]interface{}{"fake": "body"}
	attach := RequestModel{
		MapTobody: []models.MapToBody{
			{
				QueryField:  "fake",
				TargetField: "attach",
			},
		},
	}
	res := adpter.BodyMapper(body, attach)
	assert.Equal(t, res["attach"], "body")
}

func TestAdapter_BodyMapperWithInvalidQuery(t *testing.T) {
	adpter := Adapter{}
	body := map[string]interface{}{"fake": "body"}
	attach := RequestModel{
		MapTobody: []models.MapToBody{
			{
				QueryField:  "invalid",
				TargetField: "attach",
			},
		},
	}
	res := adpter.BodyMapper(body, attach)
	assert.Equal(t, res["attach"], nil)
}

func getMockSQSClient() sqsiface.SQSAPI {
	return &mocks.MockSQS{
		Messages: map[string][]*sqs.Message{},
	}
}
func TestAdapter_SendRequest(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{}}))

	var sqsOutbound ports.SecSQSPORT
	sqsOutbound = sqsOut.NewAdapter(sess, getMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	adp := NewAdapter(client, sqsOutbound, http.NewRequest)
	err := adp.SendRequest(mocks.BodyWithCallback)
	assert.Equal(t, err, nil)

}

func TestAdapter_SendRequestWithClientError(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{}}))

	var sqsOutbound ports.SecSQSPORT
	sqsOutbound = sqsOut.NewAdapter(sess, getMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClientError{
		StatusCode: 200,
	}
	adp := NewAdapter(client, sqsOutbound, http.NewRequest)
	err := adp.SendRequest(mocks.BodyWithCallback)
	assert.Equal(t, err.Error(), `I'm a do error`)

}

func TestAdapter_SendRequestWIthNoCallback(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{}}))

	var sqsOutbound ports.SecSQSPORT
	sqsOutbound = sqsOut.NewAdapter(sess, getMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	adp := NewAdapter(client, sqsOutbound, http.NewRequest)
	err := adp.SendRequest(mocks.SuccessBodyWithNoCallback)
	assert.Equal(t, err, nil)

}

func TestAdapter_SendErrorRequestWIthNoCallback(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{}}))

	var sqsOutbound ports.SecSQSPORT
	sqsOutbound = sqsOut.NewAdapter(sess, getMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	adp := NewAdapter(client, sqsOutbound, http.NewRequest)
	err := adp.SendRequest(mocks.ErrorBodyWithNoCallback)
	assert.NotEqual(t, err, nil)

}

func TestAdapter_ErrorSendRequest(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{}}))

	var sqsOutbound ports.SecSQSPORT
	sqsOutbound = sqsOut.NewAdapter(sess, getMockSQSClient(), "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	client := &mocks.MockClient{
		StatusCode: 200,
	}
	adp := NewAdapter(client, sqsOutbound, http.NewRequest)
	err := adp.SendRequest(mocks.BodyWithErrorCallback)
	assert.Equal(t, err, nil)

}
