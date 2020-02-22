package email

import "github.com/sendgrid/sendgrid-go/helpers/mail"

type template struct {
	ID      string
	From    *mail.Email
	ReplyTo *mail.Email
}

var (
	defaultFrom    = mail.NewEmail("Sam Mendelson", "sam@canvascbl.com")
	defaultReplyTo = defaultFrom
	welcome        = template{
		ID:      "d-ac67879cf72b4217b913bd768840ce3b",
		From:    defaultFrom,
		ReplyTo: defaultReplyTo,
	}
	purchaseAcknowledgement = template{
		ID:      "d-4d2ce941dc134f68945fed2aea2dfc01",
		From:    defaultFrom,
		ReplyTo: defaultReplyTo,
	}
	cancellationAcknowledgement = template{
		ID:      "d-acdc019b9ded4f71932ec5ec1dc6736f",
		From:    defaultFrom,
		ReplyTo: defaultReplyTo,
	}
	gradeChange = template{
		ID:      "d-f1313913c66440a1a3c2b48727ea20e9",
		From:    mail.NewEmail("CanvasCBL Grades", "grades@canvascbl.com"),
		ReplyTo: defaultReplyTo,
	}
)
