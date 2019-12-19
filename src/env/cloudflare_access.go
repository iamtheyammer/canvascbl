package env

var (
	CloudflareAccessAuthDomain = getEnv("CLOUDFLARE_ACCESS_AUTH_DOMAIN", "")
	CloudflareAccessAUD        = getEnv("CLOUDFLARE_ACCESS_AUD", "")
)
