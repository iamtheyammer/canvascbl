package email

import (
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func send(t template, templateData map[string]interface{}, email string, name string) {
	m := mail.NewV3Mail()
	m.SetFrom(t.From)
	m.SetTemplateID(t.ID)
	m.ReplyTo = t.ReplyTo

	p := mail.NewPersonalization()
	p.To = []*mail.Email{
		mail.NewEmail(name, email),
	}

	for k, v := range templateData {
		p.SetDynamicTemplateData(k, v)
	}

	m.AddPersonalizations(p)

	body := mail.GetRequestBody(m)

	req := sendgrid.GetRequest(env.SendGridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	req.Method = "POST"
	req.Body = body

	_, err := sendgrid.API(req)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error sending welcome email"))
	}
}
