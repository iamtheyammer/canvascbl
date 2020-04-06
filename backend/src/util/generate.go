package util

import (
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"net/url"
)

func GenerateCanvasURL(path string) url.URL {
	return url.URL{
		Host:   env.CanvasDomain,
		Scheme: "https",
		Path:   path,
	}
}
