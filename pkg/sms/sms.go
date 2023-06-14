package sms

import (
	"errors"
	"net/http"
	"net/url"
)

func (c *Client) Send(phone, message string) (err error) {
	// preparation of request params
	params := url.Values{}
	params.Add("login", c.credentials.Username)
	params.Add("psw", c.credentials.Password)
	params.Add("phones", phone)
	params.Add("mes", message)

	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	path := c.credentials.Endpoint + "/sys/send.php?" + params.Encode()

	req := Request{
		Method:  "GET",
		URL:     path,
		Body:    nil,
		Headers: headers,
	}

	// run request handler
	res, err := c.handler(req)
	if err != nil {
		return err
	}

	// check response status
	if res.Status != http.StatusOK {
		err = errors.New(string(res.Body))
	}

	return
}
