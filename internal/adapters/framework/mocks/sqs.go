package mocks

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type MockSQS struct {
	sqsiface.SQSAPI
	Messages map[string][]*sqs.Message
}

func (m *MockSQS) SendMessage(in *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	m.Messages[*in.QueueUrl] = append(m.Messages[*in.QueueUrl], &sqs.Message{
		Body: in.MessageBody,
	})
	return &sqs.SendMessageOutput{}, nil
}
func (m *MockSQS) ReceiveMessage(in *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	if len(m.Messages[*in.QueueUrl]) == 0 {
		return &sqs.ReceiveMessageOutput{}, nil
	}
	response := m.Messages[*in.QueueUrl][0:1]
	m.Messages[*in.QueueUrl] = m.Messages[*in.QueueUrl][1:]
	return &sqs.ReceiveMessageOutput{
		Messages: response,
	}, nil
}

func (m *MockSQS) DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	output := sqs.DeleteMessageOutput{}
	err := errors.New("I'm here")
	return &output, err
}

func GetMockSQSClient() sqsiface.SQSAPI {
	return &MockSQS{
		Messages: map[string][]*sqs.Message{},
	}
}
