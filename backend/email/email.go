package email

import (
	sg "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	fromName  = "smithy"
	fromEmail = "thanhanmoc1@gmail.com"
	subject   = "Account recovery code"
)

type SendGrid struct {
	apiKey string
}

func New(apiKey string) *SendGrid {
	return &SendGrid{apiKey: apiKey}
}

func (sc *SendGrid) Send(toName, toEmail, code string) error {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)
	content := mail.NewContent("text/plain", code)
	m := mail.NewV3MailInit(from, subject, to, content)

	body := mail.GetRequestBody(m)
	request := sg.GetRequest(sc.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = body

	_, err := sg.API(request)
	if err != nil {
		return err
	}

	return nil
}
