package secret

import (
	"errors"
	"net/http"
)

type Request struct {
	Key string `json:"key"`
	OTP string `json:"otp"`
}

func (s *Request) Bind(r *http.Request) error {
	if s.Key == "" {
		return errors.New("key: cannot be blank")
	}

	if s.OTP == "" {
		return errors.New("otp: cannot be blank")
	}

	return nil
}

type Response struct {
	Key string `json:"key"`
	OTP string `json:"otp,omitempty"`
}

func ParseFromEntity(data Entity) (res Response) {
	res = Response{
		Key: data.ID,
	}
	return
}
