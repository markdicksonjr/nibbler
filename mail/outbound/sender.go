package outbound

type Sender interface {
	SendMail(from *Email, subject string, to []*Email, plainTextContent string, htmlContent string) (*Response, error)
}

type Email struct {
	Name    string
	Address string
}

type Response struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
}
