package ports

type RequestPORT interface {
	SendRequest() error
}
