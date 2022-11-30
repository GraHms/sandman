package http

import (
	"bytes"
	"net/http"
)

func (r *Adapter) MakeHTTPRequest(endpoint string, method string, inputBody bytes.Reader, headers http.Header) (error, *http.Response) {

	req, err := r.NewRequest(method, endpoint, &inputBody)
	if err != nil {

		return err, nil
	}
	req.Header = headers
	response, err := r.client.Do(req)
	if err != nil {

		return err, response
	}

	return nil, response
}
