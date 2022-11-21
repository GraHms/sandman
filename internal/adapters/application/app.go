package application

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sqs"
	"serviceman/internal/pkg/models"
	"serviceman/internal/ports"
)

type Adapter struct {
	secundaryAdapter ports.RequestPORT
}

func NewAdapter(secundaryAdapter ports.RequestPORT) *Adapter {
	return &Adapter{secundaryAdapter: secundaryAdapter}
}
func (apia *Adapter) ProcessMessage(message *sqs.Message) error {
	body, err := apia.ConvertSQSBody(*message)
	if err != nil {
		return nil
	}
	err = apia.secundaryAdapter.SendRequest(body)
	if err != nil {
		return err
	}

	return err
}

func (apia *Adapter) ConvertSQSBody(msg sqs.Message) (models.Body, error) {
	var body models.Body
	err := json.Unmarshal([]byte(*msg.Body), &body)
	if err != nil {
		fmt.Println(err)
		return models.Body{}, err
	}
	//fmt.Printf("%v", body)
	return body, nil
}
