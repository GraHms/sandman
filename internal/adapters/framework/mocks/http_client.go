package mocks

import (
	"errors"
	"io"
	"net/http"
)

type MockClient struct {
	StatusCode int
	body       map[string]interface{}
}

func (m *MockClient) Do(_ *http.Request) (*http.Response, error) {
	resp := &http.Response{StatusCode: m.StatusCode, Body: MockBody(`{"foo":"bar"}`).Body}
	return resp, nil
}

type MockClientError struct {
	StatusCode int
	body       map[string]interface{}
}

func (m *MockClientError) Do(_ *http.Request) (*http.Response, error) {
	resp := &http.Response{StatusCode: m.StatusCode, Body: MockBody(`{"foo":"bar"}`).Body}
	return resp, errors.New("I'm a do error")
}

func MockNewRequestWithError(_, _ string, _ io.Reader) (*http.Request, error) {
	return &http.Request{}, errors.New("Mock Request Method Error")
}
