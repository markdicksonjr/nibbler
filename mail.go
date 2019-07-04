package nibbler

type MailSender interface {
	SendMail(from *EmailAddress, subject string, to []*EmailAddress, plainTextContent string, htmlContent string) (*MailSendResponse, error)
}

type EmailAddress struct {
	Name    string
	Address string
}

type MailSendResponse struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
}
