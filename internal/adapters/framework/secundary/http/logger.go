package http

import "encoding/json"

type Logger struct {
	Name           string `json:"name"`
	TraceId        string `json:"traceId"`
	GroupReference string `json:"groupReference"`
	Origin         string `json:"origin"`
	SandmanVersion string `json:"sandmanVersion"`
	Intent         string `json:"intent"`
	Description    string `json:"description"`
	Journey        string `json:"journey"`
	Owner          string `json:"owner"`
	Method         string `json:"method"`
	Endpoint       string `json:"endpoint"`
	Type           string `json:"type"`
	State          string `json:"state"`
	StateReason    string `json:"stateReason"`
	Action         string `json:"action"`
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) log(state, stateReason, action string, data RequestModel) {
	l.set(state, data)
	l.StateReason = stateReason
	l.Action = action
	jsonLog, _ := json.Marshal(l)
	println(string(jsonLog))

}

func (l *Logger) set(state string, data RequestModel) {
	l.Journey = data.Journey
	l.Name = data.Name
	l.TraceId = data.TraceId
	l.Owner = data.Owner
	l.Intent = data.Intent
	l.Description = data.Description
	l.SandmanVersion = data.SandmanVersion
	l.Origin = data.Origin
	l.GroupReference = data.GroupReference
	l.Endpoint = data.Endpoint
	l.State = state
	l.Type = "Main-Request"
	if data.IsCallback {
		l.Type = "Callback"
	}
}
