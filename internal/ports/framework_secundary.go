package ports

import "serviceman/internal/pkg/models"

type RequestPORT interface {
	SendRequest(body models.Body) error
}
