package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/assert/v2"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/mocks"
	"net/http"
	"testing"
)

func TestAdapter_MakeHTTPRequest(t *testing.T) {
	adp := Adapter{
		client: &mocks.MockClient{
			StatusCode: 200,
		},
		NewRequest: http.NewRequest,
	}
	bodyBytes, _ := json.Marshal(`{"fake":"body"}`)

	bodyReader := bytes.NewReader(bodyBytes)
	err, resp := adp.MakeHTTPRequest("hello.com", "POST", *bodyReader, http.Header{})
	if err != nil {
		println(err)
		return
	}
	assert.Equal(t, resp.StatusCode, 200)

}

func TestAdapter_MakeErrorClient(t *testing.T) {
	adp := Adapter{
		client: &mocks.MockClientError{
			StatusCode: 200,
		},
		NewRequest: http.NewRequest,
	}
	bodyBytes, _ := json.Marshal(`{"fake":"body"}`)

	bodyReader := bytes.NewReader(bodyBytes)
	err, _ := adp.MakeHTTPRequest("hello.com", "POST", *bodyReader, http.Header{})
	if err != nil {
		println(err)
		return
	}
	assert.Equal(t, err.Error(), `I'm a do error`)

}

func TestAdapter_MakeHTTPRequestWithWrongMethod(t *testing.T) {
	adp := Adapter{
		client: &mocks.MockClient{
			StatusCode: 200,
		},
		NewRequest: http.NewRequest,
	}
	bodyBytes, _ := json.Marshal(`{"fake":"body"}`)

	bodyReader := bytes.NewReader(bodyBytes)
	err, resp := adp.MakeHTTPRequest("hello.com", "NONONON", *bodyReader, http.Header{})
	if err != nil {
		println(err)
		return
	}
	assert.Equal(t, resp.StatusCode, 200)

}

func TestAdapter_MakeHTTPRequestWithNewRequestError(t *testing.T) {
	adp := Adapter{
		client: &mocks.MockClient{
			StatusCode: 200,
		},
		NewRequest: mocks.MockNewRequestWithError,
	}
	bodyBytes, _ := json.Marshal(`{"fake":"body"}`)

	bodyReader := bytes.NewReader(bodyBytes)
	err, _ := adp.MakeHTTPRequest("hello.com", "NONONON", *bodyReader, http.Header{})

	assert.Equal(t, err.Error(), `Mock Request Method Error`)

}
