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
	Name           string   `json:"name"`
	TraceId        string   `json:"traceId"`
	GroupReference string   `json:"groupReference"`
	Origin         string   `json:"origin"`
	SandmanVersion string   `json:"sandmanVersion"`
	Intent         string   `json:"intent"`
	Description    string   `json:"description"`
	Journey        string   `json:"journey"`
	Owner          string   `json:"owner"`
	Request        Request  `json:"request"`
	Response       Response `json:"response"`
}

type Callbacks struct {
	IsDefault         bool   `json:"isDefault"`
	Name              string `json:"name"`
	WhenRequestStatus int    `json:"whenRequestStatus,omitempty"`
	SuccessStatus     int    `json:"successStatus"`
	Request

	MapToBody []MapToBody `json:"mapToBody"`
}

type MapToBody struct {
	QueryField  string `json:"queryField"`
	TargetField string `json:"targetField"`
}
