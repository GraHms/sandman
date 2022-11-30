package mocks

import "lab.dev.vm.co.mz/compse/sandman/internal/pkg/models"

var BodyWithCallback = models.Body{Response: models.Response{
	SuccessStatus: 200,
	HasCallBack:   true,
	Callbacks: []models.Callbacks{
		{
			WhenRequestStatus: 200,
			Request: models.Request{
				Body:    map[string]interface{}{"foo": "bar"},
				Sub:     "grahms.com",
				Method:  "POST",
				Retries: 1,
			},
		},
	},
},
	Request: models.Request{
		Body:        map[string]interface{}{"foo": "bar"},
		Sub:         "hello.com",
		Method:      "POST",
		Retries:     1,
		ContentType: "application/json",
	},
}

var SuccessBodyWithNoCallback = models.Body{Response: models.Response{
	SuccessStatus: 200,
	HasCallBack:   false,
},
	Request: models.Request{
		Body:        map[string]interface{}{"foo": "bar"},
		Sub:         "hello.com",
		Method:      "POST",
		Retries:     1,
		ContentType: "application/json",
	},
}

var ErrorBodyWithNoCallback = models.Body{Response: models.Response{
	SuccessStatus: 400,
	HasCallBack:   false,
},
	Request: models.Request{
		Body:        map[string]interface{}{"foo": "bar"},
		Sub:         "hello.com",
		Method:      "POST",
		Retries:     1,
		ContentType: "application/json",
	},
}

var BodyWithErrorCallback = models.Body{Response: models.Response{
	SuccessStatus: 201,
	HasCallBack:   true,
	Callbacks: []models.Callbacks{
		{
			WhenRequestStatus: 200,
			Request: models.Request{
				Body:    map[string]interface{}{"foo": "bar"},
				Sub:     "grahms.com",
				Method:  "POST",
				Retries: 1,
			},
		},
		{
			WhenRequestStatus: 200,
			Request: models.Request{
				Body:    map[string]interface{}{"foo": "bar"},
				Sub:     "grahms.com",
				Method:  "POST",
				Retries: 1,
			},
		},
	},
},
	Request: models.Request{
		Body:        map[string]interface{}{"foo": "bar"},
		Sub:         "hello.com",
		Method:      "POST",
		Retries:     1,
		ContentType: "application/json",
	},
}
