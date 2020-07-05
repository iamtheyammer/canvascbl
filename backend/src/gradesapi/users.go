package gradesapi

import (
	"errors"
	"fmt"
	userssvc "github.com/iamtheyammer/canvascbl/backend/src/db/services/users"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type userProfile struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	LTIUserID    string `json:"lti_user_id,omitempty"`
	CanvasUserID uint64 `json:"canvas_user_id"`
	Status       int    `json:"status,omitempty"`
}

type profileHandlerResponse struct {
	Profile userProfile `json:"profile"`
}

//UserProfileHandler handles /api/v1/users/:userID/profile
func UserProfileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	requestedUserID := ps.ByName("userID")
	if len(requestedUserID) < 1 {
		util.SendBadRequest(w, "missing userID as url param")
		return
	}

	userID, rdP, sess, errCtx := authorizer(w, r, []oauth2.Scope{oauth2.ScopeProfile}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "users/:userID/profile",
		Query:     &r.URL.RawQuery,
	})
	if (userID == nil || rdP == nil || errCtx == nil) && sess == nil {
		return
	}

	errCtx.AddCustomField("requested_user_id", requestedUserID)

	if requestedUserID != "self" && fmt.Sprintf("%d", *userID) != requestedUserID {
		util.SendUnauthorized(w, "insufficient permissions")
		return
	}

	profiles, err := userssvc.List(db, &userssvc.ListRequest{ID: *userID})
	if err != nil {
		handleISE(w, errCtx.Apply(fmt.Errorf("error listing profiles in user profile handler: %w", err)))
		return
	}

	if len(*profiles) < 1 {
		handleISE(w, errCtx.Apply(errors.New("zero users returned when listing profiles in user profile handler")))
		return
	}

	p := (*profiles)[0]

	up := userProfile{
		ID:           p.ID,
		Name:         p.Name,
		Email:        p.Email,
		CanvasUserID: p.CanvasUserID,
	}

	if r.URL.Query().Get("include[]") == "extra_info" {
		up.LTIUserID = p.LTIUserID
		up.Status = p.Status
	}

	sendJSON(w, &profileHandlerResponse{Profile: up})
}
