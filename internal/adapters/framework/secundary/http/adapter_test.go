package http

import (
	"net/http"
	"testing"
)

type ClientMock struct {
}

func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}
func TestNewAdapter(t *testing.T) {
	//NewAdapter(ClientMock{})
}
