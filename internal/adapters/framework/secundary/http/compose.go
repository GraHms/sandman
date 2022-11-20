package http

import "serviceman/internal/pkg/models"

type Compose struct {
	List *List
	body models.Body
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
	reqList := NewList(*mainReq)
	c.AttachCallBacks()
	c.List = reqList
}

func (c *Compose) AttachCallBacks() {
	if c.body.Response.HasCallBack {
		for _, callback := range c.body.Response.Callbacks {
			c.List.Insert(*c.PrepareCallback(callback))
		}
	}
}

func (c *Compose) PrepareMain() *RequestModel {
	var req RequestModel
	req.IsCallback = false
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
	req.CallBackExecuteWhenStatusIs = c.body.Response.SuccessStatus
	req.IsCallback = true
	req.ExpectSuccessStatus = body.SuccessStatus
	req.Retries = body.Retries
	req.Endpoint = body.Sub
	req.Method = body.Method
	req.Body = body.Body
	req.ContentType = body.ContentType
	req.Headers = body.Headers
	return &req
}
