package http

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_RequestError(t *testing.T) {
	reqErr := RequestError{StatusCode: 200,
		Err: errors.New("I'm an error")}
	assert.Equal(t, reqErr.Error(), "I'm an error")
}
