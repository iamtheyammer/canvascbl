package db

import (
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
)

func VerifySession(sessionString string) (*sessions.VerifiedSession, error) {
	return sessions.Verify(util.DB, sessionString)
}
