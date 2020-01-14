package env

var (
	SendGridAPIKey                                = getEnvOrPanic("SENDGRID_API_KEY")
	SendGridFromName                              = getEnvOrPanic("SENDGRID_FROM_NAME")
	SendGridFromEmail                             = getEnvOrPanic("SENDGRID_FROM_EMAIL")
	SendGridWelcomeTemplateID                     = getEnvOrPanic("SENDGRID_WELCOME_TEMPLATE_ID")
	SendGridPurchaseAcknowledgementTemplateID     = getEnvOrPanic("SENDGRID_PURCHASE_ACKNOWLEDGEMENT_TEMPLATE_ID")
	SendGridCancellationAcknowledgementTemplateID = getEnvOrPanic("SENDGRID_CANCELLATION_ACKNOWLEDGEMENT_TEMPLATE_ID")
)
