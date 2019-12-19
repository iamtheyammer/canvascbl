package util

import (
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/cfjwt"
)

var CloudflareAccessVerifier = cfjwt.Verifier{
	PolicyAUD:  env.CloudflareAccessAUD,
	AuthDomain: env.CloudflareAccessAuthDomain,
}.Init()
