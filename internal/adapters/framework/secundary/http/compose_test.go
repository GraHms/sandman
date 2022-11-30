package http

import (
	"github.com/go-playground/assert/v2"
	"lab.dev.vm.co.mz/compse/sandman/internal/adapters/framework/mocks"
	"testing"
)

func TestNewComposer(t *testing.T) {

	composer := NewComposer(mocks.BodyWithCallback)
	candidates := composer.Candidates
	head := candidates.head

	assert.Equal(t, head.Method, "POST")
	assert.Equal(t, head.IsCallback, false)
	callbacks := candidates.callbacks
	assert.Equal(t, callbacks[200][0].IsCallback, true)
	assert.Equal(t, callbacks[200][0].Retries, 1)
}
