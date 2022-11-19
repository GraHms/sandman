package http

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
)

func (r *Adapter) SetHttpClient() {
	ssl := &tls.Config{
		InsecureSkipVerify: true,
	}
	r.client = &http.Client{Timeout: 120 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: ssl,
		}}
}

func (r *Adapter) MakeHTTPRequest(endpoint string, method string, inputBody bytes.Reader, headers http.Header) (error, *http.Response) {

	req, err := http.NewRequest(method, endpoint, &inputBody)
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
