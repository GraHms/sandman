package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"serviceman/internal/pkg/models"
)

type Adapter struct {
	client *http.Client
}

func NewAdapter(client *http.Client) *Adapter {
	adp := Adapter{
		client: client,
	}
	adp.SetHttpClient()
	return &adp
}

func (seca *Adapter) SendRequest(body models.Body) error {
	println("making an http call")
	lineUp := NewComposer(body)
	elected := lineUp.List.GetHead()
	var headExpectedStatus int
	var headResponse *http.Response
	defer bodyCloser(headResponse.Body)
	for elected != nil {
		if elected.ReqModel.IsCallback != true {
			headExpectedStatus = elected.ReqModel.ExpectSuccessStatus
		}
		if (elected.ReqModel.IsCallback == true && headExpectedStatus != elected.ReqModel.CallBackExecuteWhenStatusIs) || (headResponse.StatusCode != headExpectedStatus) {
			// if callback isn't triggered by head status code response
			lineUp.List.PopHead()
			elected = elected.next
			continue
		}

		respBody, err := seca.Execute(&elected.ReqModel)
		if elected.ReqModel.IsCallback != true {
			headResponse = respBody
		} else {
			bodyCloser(respBody.Body)
		}
		if err != nil && elected.ReqModel.IsCallback == false {
			// if error on head request, break entire process loop
			return err
		}
		//	send elected to queue and continue execution
		if err != nil {

			println("This callback failed, we are sending it to the queue")
		}
		lineUp.List.PopHead()
		elected = elected.next
	}
	fmt.Println()
	return nil
}

func (seca *Adapter) Execute(reqModel *RequestModel) (*http.Response, error) {

	//	first step, call main request
	headers := http.Header{
		"Content-Type": {reqModel.ContentType},
		"charset":      {"utf-8"},
	}
	var resp *http.Response
	bodyInter := reqModel.Body.(map[string]interface{})
	bodyBytes, err := json.Marshal(bodyInter)

	if err != nil {
		fmt.Printf("couldn't marshal body: ", bodyBytes)
		return resp, err
	}
	bodyReader := bytes.NewReader(bodyBytes)
	for i := 0; i < reqModel.Retries; i++ {
		println("calling ", reqModel.Endpoint, " on attempt number:", i)

		err, resp = seca.MakeHTTPRequest(reqModel.Endpoint, reqModel.Method, *bodyReader, headers)

		if err != nil {
			return resp, err
		}
		if resp.StatusCode == reqModel.ExpectSuccessStatus {
			println("request successfully made to ", reqModel.Endpoint, " The status code is: ", resp.StatusCode)
			return resp, nil
		}
	}

	e := "couldn't call " + reqModel.Endpoint + " after " + string(rune(reqModel.Retries)) + " attempts. expected status is:" + string(rune(reqModel.ExpectSuccessStatus)) + " and given is " + string(rune(resp.StatusCode))
	return resp, errors.New(e)
}

func bodyCloser(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		return
	}
}
func (seca *Adapter) PrepareMainRequest(body models.Body) *RequestModel {
	var req RequestModel
	req.ExpectSuccessStatus = body.Response.SuccessStatus
	req.Retries = body.Request.Retries
	req.Endpoint = body.Request.Sub
	req.Method = body.Request.Method
	req.Body = body.Request.Body
	req.ContentType = body.Request.ContentType
	req.Headers = body.Request.Headers
	return &req

}

type RequestModel struct {
	IsCallback                  bool
	CallBackExecuteWhenStatusIs int
	Retries                     int
	ExpectSuccessStatus         int
	Method                      string
	ContentType                 string
	Endpoint                    string
	Body                        interface{}
	Headers                     string
}
