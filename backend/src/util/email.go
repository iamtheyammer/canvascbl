package util

import (
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var SendGridFrom = mail.NewEmail(env.SendGridFromName, env.SendGridFromEmail)
