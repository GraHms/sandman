package models

type Request struct {
	Retries     int         `json:"retries"`
	Method      string      `json:"method"`
	ContentType string      `json:"contentType"`
	Body        interface{} `json:"body"`
	Sub         string      `json:"sub"`
	Headers     string      `json:"headers"`
}

type Response struct {
	SuccessStatus int         `json:"successStatus"`
	HasCallBack   bool        `json:"hasCallBack"`
	Callbacks     []Callbacks `json:"callbacks"`
}

type Body struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

type Callbacks struct {
	IsDefault         bool   `json:"isDefault"`
	Name              string `json:"name"`
	WhenRequestStatus int    `json:"whenRequestStatus,omitempty"`
	SuccessStatus     int    `json:"successStatus"`
	Request
}
