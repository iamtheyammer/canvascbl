package canvasapis

import (
	"encoding/json"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/canvasapis/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	usersdecoder "github.com/iamtheyammer/canvascbl/backend/src/db/canvas/users"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/canvas_tokens"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/google_users"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/sessions"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/env"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

type canvasToken struct {
	ID         uint64    `json:"id"`
	UserID     uint64    `json:"userId"`
	Token      string    `json:"token"`
	Subdomain  string    `json:"subdomain"`
	InsertedAt time.Time `json:"insertedAt"`
}

type insertCanvasTokenHandlerBody struct {
	Token string `json:"token"`
	// epoch
	ExpiresAt uint64 `json:"expiresAt"`
}

func InsertCanvasTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := middlewares.Session(w, r)
	if sess == nil {
		return
	}

	var body insertCanvasTokenHandlerBody
	err := middlewares.DecodeJSONBody(r.Body, &body)
	if err != nil {
		util.SendBadRequest(w, "invalid JSON body (do you have one?)")
		return
	}

	if len(body.Token) < 1 {
		util.SendBadRequest(w, "missing token field in JSON POST body")
		return
	}

	if body.ExpiresAt > 0 && time.Unix(int64(body.ExpiresAt), 0).Before(time.Now()) {
		util.SendBadRequest(w, "token is already expired (expiresAt should be secs from epoch)")
	}

	if !util.ValidateCanvasToken(body.Token) {
		util.SendBadRequest(w, "that doesn't seem like a canvas token")
		return
	}

	res, profileJSON, err := userssvc.GetSelfProfile(&util.RequestDetails{
		Token:     body.Token,
		Subdomain: env.CanvasOAuth2Subdomain,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error validating token by fetching user profile"))
		util.SendInternalServerError(w)
		return
	}

	if res.StatusCode != http.StatusOK {
		util.SendBadRequest(w, "invalid canvas token")
		return
	}

	p, err := usersdecoder.ProfileFromJSON(&profileJSON)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error getting canvas profile from json"))
		util.SendInternalServerError(w)
		return
	}

	db.UpsertProfile(&profileJSON)

	gUsersP, err := db.ListGoogleProfiles(&google_users.ListRequest{
		Email: strings.ToLower(p.PrimaryEmail),
	})

	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing google users"))
		util.SendInternalServerError(w)
		return
	}

	usP, err := db.ListUsers(&users.ListRequest{
		Email: p.PrimaryEmail,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing users"))
		util.SendInternalServerError(w)
		return
	}

	gUsers := *gUsersP
	us := *usP
	if len(gUsers) != 1 {
		util.SendBadRequest(w, "seems like you haven't signed up (?)")
		return
	}

	if len(us) != 1 {
		util.HandleError(errors.New("more or less than 1 profile was returned in insert canvas token handler"))
		util.SendInternalServerError(w)
		return
	}

	gUser := gUsers[0]

	ir := canvas_tokens.InsertRequest{
		GoogleUsersID: gUser.ID,
		CanvasUserID:  &us[0].CanvasUserID,
		UserID:        &us[0].ID,
		Token:         body.Token,
	}

	if strings.ToLower(us[0].Email) != strings.ToLower(p.PrimaryEmail) ||
		strings.ToLower(us[0].Email) != strings.ToLower(gUser.Email) {
		util.SendBadRequest(w, "preexisting user email does not match "+
			"google email or canvas email")
		return
	}

	if body.ExpiresAt > 0 {
		t := time.Unix(int64(body.ExpiresAt), 0)
		ir.ExpiresAt = &t
	}

	err = db.InsertCanvasToken(&ir)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error inserting canvas token"))
		util.SendInternalServerError(w)
		return
	}

	err = db.UpdateGoogleProfile(&google_users.UpdateRequest{
		UsersID: us[0].ID,
		WhereID: gUser.ID,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error updating google profile"))
		util.SendInternalServerError(w)
		return
	}

	err = db.UpdateSession(&sessions.UpdateRequest{
		CanvasUserID:       uint64(p.ID),
		WhereSessionString: sess.SessionString,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error updating session"))
		util.SendInternalServerError(w)
		return
	}

	util.SendNoContent(w)
	return
}

func GetCanvasTokensHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess := middlewares.Session(w, r)
	if sess == nil {
		return
	}

	tokens, err := db.ListCanvasTokens(&canvas_tokens.ListRequest{
		UserID:         sess.UserID,
		CanvasUserID:   sess.CanvasUserID,
		NotExpiredOnly: true,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error listing canvas tokens"))
		util.SendInternalServerError(w)
		return
	}

	var toks []canvasToken

	for _, t := range *tokens {
		toks = append(toks, canvasToken{
			ID:         t.ID,
			UserID:     t.UserID,
			Token:      t.Token,
			Subdomain:  env.CanvasOAuth2Subdomain,
			InsertedAt: t.InsertedAt,
		})
	}

	if len(toks) < 1 {
		util.SendJSONResponse(w, []byte("[]"))
		return
	}

	jToks, err := json.Marshal(toks)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling tokens from get canvas tokens"))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jToks)
	return
}

func DeleteCanvasTokenHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ok, rd := util.GetRequestDetailsFromRequest(r)
	if !ok {
		util.SendUnauthorized(w, util.RequestDetailsFailedValidationMessage)
		return
	}

	sess := middlewares.Session(w, r)
	if sess == nil {
		return
	}

	if sess.UserID < 1 {
		util.SendUnauthorized(w, "appears that you haven't logged in yet")
		return
	}

	err := db.DeleteCanvasToken(&canvas_tokens.DeleteRequest{Token: rd.Token})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error deleting canvas token"))
		util.SendInternalServerError(w)
		return
	}

	util.SendNoContent(w)
	return
}
