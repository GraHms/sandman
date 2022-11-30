package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/go-playground/assert/v2"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/mocks"
	"lab.dev.vm.co.mz/compse/sandman/internal/pkg/models"
	"lab.dev.vm.co.mz/compse/sandman/internal/ports"
	"testing"
)

func getMockSQSClient() sqsiface.SQSAPI {
	return &mocks.MockSQS{
		Messages: map[string][]*sqs.Message{},
	}
}
func TestAdapter_SendMessage(t *testing.T) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{}}))
	q := getMockSQSClient()
	var sqsOutbound ports.SecSQSPORT
	sqsOutbound = NewAdapter(sess, q, "https://queue.amazonaws.com/80398EXAMPLE/MyQueue")
	body := models.Body{}
	err := sqsOutbound.SendMessage(body)
	if err != nil {
		return
	}
	assert.Equal(t, err, nil)
}

func TestAdapter_BodyComposer(t *testing.T) {
	a := Adapter{}
	expected := `{"name":"","traceId":"","groupReference":"","origin":"","sandmanVersion":"","intent":"","description":"","journey":"","owner":"","request":{"retries":0,"method":"","contentType":"","body":null,"sub":"","headers":""},"response":{"successStatus":0,"hasCallBack":false,"callbacks":null}}`
	composer, _ := a.BodyComposer(models.Body{})
	assert.Equal(t, composer, expected)

}
