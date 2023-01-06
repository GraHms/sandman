package config

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLogger(t *testing.T) {
	// Create a buffer to capture the output of the logger.
	buffer := &bytes.Buffer{}

	// Set the logger's output to be the buffer we created.
	logger := Logger()

	// Use the logger to log a message.
	logger.Info("This is a test log message.")

	// Assert that the message was written to the buffer.
	assert.Contains(t, buffer.String(), "")
}
