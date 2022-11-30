package http

import "lab.dev.vm.co.mz/compse/sandman/internal/pkg/models"

type Compose struct {
	body       models.Body
	Candidates ListCandidates
}

type ListCandidates struct {
	head      RequestModel
	callbacks map[int][]RequestModel
}

func NewComposer(body models.Body) *Compose {
	compose := Compose{
		body: body,
	}
	compose.Compose()
	return &compose
}

func (c *Compose) Compose() {

	mainReq := c.PrepareMain()
	c.Candidates.head = *mainReq
	if c.body.Response.HasCallBack {
		c.AttachCallBacks()
	}
}

func (c *Compose) AttachCallBacks() {
	callbacks := make(map[int][]RequestModel)
	for _, callback := range c.body.Response.Callbacks {
		triggerStatus := callback.WhenRequestStatus
		val := c.PrepareCallback(callback)
		if _, ok := callbacks[triggerStatus]; !ok {
			callbacks[triggerStatus] = []RequestModel{
				*val,
			}
			continue
		}
		callbacks[triggerStatus] = append(callbacks[triggerStatus], *val)
	}

	c.Candidates.callbacks = callbacks
}

func (c *Compose) PrepareMain() *RequestModel {
	var req RequestModel
	req.IsCallback = false
	req.Intent = c.body.Intent
	req.TraceId = c.body.TraceId
	req.Description = c.body.Description
	req.GroupReference = c.body.GroupReference
	req.Owner = c.body.Owner
	req.SandmanVersion = c.body.SandmanVersion
	req.Name = c.body.Name
	req.Journey = c.body.Journey
	req.ExpectSuccessStatus = c.body.Response.SuccessStatus
	req.Retries = c.body.Request.Retries
	req.Endpoint = c.body.Request.Sub
	req.Method = c.body.Request.Method
	req.Body = c.body.Request.Body
	req.ContentType = c.body.Request.ContentType
	req.Headers = c.body.Request.Headers
	return &req
}

func (c *Compose) PrepareCallback(body models.Callbacks) *RequestModel {
	var req RequestModel
	req.Intent = c.body.Intent
	req.TraceId = c.body.TraceId
	req.Description = c.body.Description
	req.GroupReference = c.body.GroupReference
	req.Owner = c.body.Owner
	req.SandmanVersion = c.body.SandmanVersion
	req.Name = body.Name
	req.Journey = c.body.Journey
	req.CallBackExecuteWhenStatusIs = c.body.Response.SuccessStatus
	req.IsCallback = true
	req.ExpectSuccessStatus = body.SuccessStatus
	req.Retries = body.Retries
	req.Endpoint = body.Sub
	req.Method = body.Method
	req.Body = body.Body
	req.ContentType = body.ContentType
	req.Headers = body.Headers
	req.MapTobody = body.MapToBody
	return &req
}
