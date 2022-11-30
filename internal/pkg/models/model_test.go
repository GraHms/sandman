package models

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestModel(t *testing.T) {
	model := Body{Name: "Ismael GraHms"}
	assert.Equal(t, model.Name, "Ismael GraHms")
}
