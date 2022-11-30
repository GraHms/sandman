package ports

import (
	"lab.dev.vm.co.mz/compse/sandman/internal/pkg/models"
	"net/http"
)

type RequestPORT interface {
	SendRequest(body models.Body) error
}

type SecSQSPORT interface {
	SendMessage(body models.Body) error
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
