package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lab.dev.vm.co.mz/compse/sandman/internal/pkg/models"
	"lab.dev.vm.co.mz/compse/sandman/internal/ports"
	"net/http"
	"sync"
)

type Adapter struct {
	client      ports.HTTPClient
	sqsOutbound ports.SecSQSPORT
	NewRequest  func(method, url string, body io.Reader) (*http.Request, error)
}

var logger = NewLogger()

func NewAdapter(client ports.HTTPClient, sqsOutbound ports.SecSQSPORT, NewRequest func(method, url string, body io.Reader) (*http.Request, error)) *Adapter {
	adp := Adapter{
		client:      client,
		sqsOutbound: sqsOutbound,
		NewRequest:  NewRequest,
	}

	return &adp
}

func (seca *Adapter) ConvertBodyResponse(resp *http.Response) map[string]interface{} {
	data := make(map[string]interface{})
	err := json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		println("could not unmarshal response body", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	return data
}

func (seca *Adapter) BodyMapper(body map[string]interface{}, model RequestModel) map[string]interface{} {
	data := make(map[string]interface{})
	for _, val := range model.MapTobody {
		queryField := val.QueryField
		targetField := val.TargetField
		if _, ok := body[queryField]; !ok {
			// could not find field to map
			println("could not map <", targetField, "> on callback because we could not find the field <", queryField, "> in main response")
			continue
		}
		println("mapping <", targetField, "> to <", queryField, "> on callback: ", model.Endpoint)
		println("")
		data[targetField] = body[queryField]

	}
	return data
}
func (seca *Adapter) SendRequest(body models.Body) error {

	lineUp := NewComposer(body)
	electedHead := lineUp.Candidates.head
	callbacks := lineUp.Candidates.callbacks

	reqBody := electedHead.Body.(map[string]interface{})
	headRespBody, err := seca.Execute(&electedHead, reqBody)
	if err != nil {
		// if is not an api error
		// send request back to queue
		if _, ok := err.(*RequestError); !ok {
			logger.log("failed", "", "keepInQueue", electedHead)
			return err
		}
		// check if error status code has callbacks to execute
		if val, ok := callbacks[headRespBody.StatusCode]; ok {
			logger.log("failed", "", "processFailEvents", electedHead)
			// if it has any failure callback, then it should delete main request
			mapHeadBody := seca.ConvertBodyResponse(headRespBody)

			seca.ExecCallbacks(val, mapHeadBody)
			logger.log("failed", "", "popdAndProcessedFailedEvents", electedHead)
			return nil
		}
		// this error has no callback
		logger.log("failed", "", "keepInQueue", electedHead)
		return err
	}
	// head processed as expected, so the head should be removed from queue
	logger.log("processed", "This event was successfully processed", "popdFromQueue", electedHead)
	// execute success callbacks
	if val, ok := callbacks[headRespBody.StatusCode]; ok {
		seca.ExecCallbacks(val, seca.ConvertBodyResponse(headRespBody))
	}

	fmt.Println()
	return nil
}
func (seca *Adapter) SQSOutBoundBodyMapper(data RequestModel, body map[string]interface{}) models.Body {
	jsonBody, _ := json.Marshal(body)
	outBound := models.Body{
		Name:           data.Name,
		TraceId:        data.TraceId,
		GroupReference: data.GroupReference,
		Intent:         data.Intent,
		Description:    data.Description,
		Journey:        "FailedAndPopped",
		Owner:          data.Owner,
		SandmanVersion: data.SandmanVersion,
		Request: models.Request{
			Sub:         data.Endpoint,
			ContentType: data.ContentType,
			Headers:     data.Headers,
			Retries:     data.Retries,
			Body:        string(jsonBody),
		},
		Response: models.Response{
			HasCallBack:   false,
			SuccessStatus: data.ExpectSuccessStatus,
		},
	}
	return outBound
}
func (seca *Adapter) ExecCallbacks(callbacks []RequestModel, respBody map[string]interface{}) {

	wg := sync.WaitGroup{}
	for _, callback := range callbacks {

		localCallback := callback
		wg.Add(1)
		go func() {

			defer wg.Done()

			attachToBody := seca.BodyMapper(respBody, localCallback)
			reqBody := localCallback.Body.(map[string]interface{})
			for k, v := range attachToBody {
				reqBody[k] = v
			}
			res, err := seca.Execute(&localCallback, reqBody)
			println(res.StatusCode)
			if err != nil {

				//if error, add this callback to the queue, as a new main
				_ = seca.sqsOutbound.SendMessage(seca.SQSOutBoundBodyMapper(localCallback, reqBody))
				logger.log("failed", "could not execute callback", "electedNewEvent", localCallback)

				return
			}
			logger.log("processed", "successfully processed callback", "popdEvent", localCallback)
		}()

	}
	wg.Wait()
}

func (seca *Adapter) Execute(reqModel *RequestModel, body map[string]interface{}) (*http.Response, error) {

	//	first step, call main request
	headers := http.Header{
		"Content-Type": {reqModel.ContentType},
		"charset":      {"utf-8"},
	}
	var resp *http.Response

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		fmt.Printf("couldn't marshal body: %v", bodyBytes)
		return resp, err
	}
	bodyReader := bytes.NewReader(bodyBytes)
	for i := 0; i < reqModel.Retries; i++ {
		println("calling ", reqModel.Endpoint, " on attempt number:", i)

		err, resp = seca.MakeHTTPRequest(reqModel.Endpoint, reqModel.Method, *bodyReader, headers)

		if err != nil {
			println("could not execute http call: ", err)
			return resp, err
		}
		if resp.StatusCode == reqModel.ExpectSuccessStatus {
			println("request successfully made to ", reqModel.Endpoint, " The status code is: ", resp.StatusCode)
			return resp, nil
		}
	}
	println("Failed http execution: expected status", reqModel.ExpectSuccessStatus, "found: ", resp.StatusCode)
	e := "couldn't call " + reqModel.Endpoint + " after " + string(rune(reqModel.Retries)) + " attempts. expected status is:" + string(rune(reqModel.ExpectSuccessStatus)) + " and given is " + string(rune(resp.StatusCode))
	println(e)
	reqErr := &RequestError{
		StatusCode: resp.StatusCode,
		Err:        errors.New(e),
	}
	return resp, reqErr
}

type RequestModel struct {
	Name                        string
	TraceId                     string
	GroupReference              string
	Origin                      string
	SandmanVersion              string
	Intent                      string
	Description                 string
	Journey                     string
	Owner                       string
	IsCallback                  bool
	CallBackExecuteWhenStatusIs int
	Retries                     int
	ExpectSuccessStatus         int
	Method                      string
	ContentType                 string
	Endpoint                    string
	Body                        interface{}
	Headers                     string
	MapTobody                   []models.MapToBody
}
