package env

var (
	StripeAPIKey             = getEnvOrPanic("STRIPE_API_KEY")
	StripeWebhookSecret      = getEnvOrPanic("STRIPE_WEBHOOK_SECRET")
	StripeCancelPurchaseURL  = getEnvOrPanic("STRIPE_CANCEL_PURCHASE_URL")
	StripePurchaseSuccessURL = getEnvOrPanic("STRIPE_PURCHASE_SUCCESS_URL")
)
