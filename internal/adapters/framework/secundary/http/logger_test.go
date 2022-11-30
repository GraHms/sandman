package http

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func Test_Logger(t *testing.T) {
	logger := Logger{}
	logger.set("test", RequestModel{})
	logger.log("error", "running tests", "unitets", RequestModel{IsCallback: true})
	assert.Equal(t, logger.State, "error")
}
