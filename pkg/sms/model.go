package sms

type Credentials struct {
	Endpoint string
	Username string
	Password string
}

type Request struct {
	Method  string
	URL     string
	Body    []byte
	Headers map[string]string
}

type Response struct {
	Status int
	Body   []byte
}
