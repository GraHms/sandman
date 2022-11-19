package application

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"serviceman/internal/pkg/models"
)

type Adapter struct {
}

func NewAdapter() *Adapter {
	return &Adapter{}
}
func (apia *Adapter) ProcessMessage(message *sqs.Message) error {
	body, err := apia.ConvertSQSBody(*message)
	if err != nil {
		return err
	}
	fmt.Printf("%v", body)
	return err
}

func (apia *Adapter) ConvertSQSBody(msg sqs.Message) (models.Body, error) {
	var body models.Body
	err := json.Unmarshal([]byte(*msg.Body), &body)
	if err != nil {
		return models.Body{}, err
	}
	fmt.Printf("%v", body)
	return body, nil
}
