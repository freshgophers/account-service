package sms

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client      *http.Client
	credentials Credentials
}

func New(c Credentials) *Client {
	client := http.DefaultClient
	client.Timeout = 30 * time.Second

	return &Client{
		client:      client,
		credentials: c,
	}
}

func (c *Client) handler(req Request) (res Response, err error) {
	// create new request
	request, err := http.NewRequest(req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return
	}

	// setup request header
	for key, value := range req.Headers {
		request.Header.Add(key, value)
	}

	// send request
	response, err := c.client.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	// read response status and body
	res.Status = response.StatusCode
	res.Body, err = io.ReadAll(response.Body)

	return
}
