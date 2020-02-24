package gradesapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/notifications"
	"github.com/iamtheyammer/canvascbl/backend/src/oauth2"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"sync"
)

type notificationType struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	Description string `json:"description,omitempty"`
}

type notificationSetting struct {
	Type   uint64 `json:"notification_type_id"`
	Medium string `json:"medium"`
}

type listNotificationSettingsResponse struct {
	NotificationSettings []notificationSetting `json:"notification_settings"`
	NotificationTypes    []notificationType    `json:"notification_types,omitempty"`
}

// ListNotificationTypesHandler lists notification types.
func ListNotificationTypesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeNotifications}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "notifications/types",
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	ts, err := notifications.ListTypes(db, &notifications.ListTypesRequest{})
	if err != nil {
		handleISE(w, fmt.Errorf("error listing notification types: %w", err))
		return
	}

	var nts []notificationType
	for _, t := range *ts {
		nts = append(nts, notificationType{
			ID:          t.ID,
			Name:        t.Name,
			ShortName:   t.ShortName,
			Description: t.Description,
		})
	}

	jRet, err := json.Marshal(&nts)
	if err != nil {
		handleISE(w, fmt.Errorf("error marshaling notification types in list notification types handler: %w", err))
		return
	}

	util.SendJSONResponse(w, jRet)
	return
}

// ListNotificationSettingsHandler lists notification settings.
func ListNotificationSettingsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeNotifications}, &oauth2.AuthorizerAPICall{
		Method:    "GET",
		RoutePath: "notifications/types",
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	includeTypes := r.URL.Query().Get("include[]") == "notification_types"

	var (
		wg       = sync.WaitGroup{}
		settings []notificationSetting
		types    []notificationType
		err      error
	)

	// settings
	wg.Add(1)
	go func() {
		defer wg.Done()

		ss, e := notifications.ListSettings(db, &notifications.ListSettingsRequest{
			UserID: *userID,
		})
		if e != nil {
			err = fmt.Errorf("error listing notification settings: %w", err)
			return
		}

		for _, s := range *ss {
			settings = append(settings, notificationSetting{
				Type:   s.Type,
				Medium: string(s.Medium),
			})
		}
	}()

	// types
	if includeTypes {
		wg.Add(1)
		go func() {
			defer wg.Done()

			ts, e := notifications.ListTypes(db, &notifications.ListTypesRequest{})
			if e != nil {
				err = fmt.Errorf("error listing notification types: %w", err)
				return
			}

			for _, t := range *ts {
				types = append(types, notificationType{
					ID:          t.ID,
					Name:        t.Name,
					ShortName:   t.ShortName,
					Description: t.Description,
				})
			}
		}()
	}

	wg.Wait()

	if err != nil {
		handleISE(w, fmt.Errorf("error in list notification settings handler: %w", err))
		return
	}

	jRet, err := json.Marshal(&listNotificationSettingsResponse{
		NotificationSettings: settings,
		NotificationTypes:    types,
	})
	if err != nil {
		handleISE(w, fmt.Errorf("error marshaling list notification settings response: %w", err))
		return
	}

	util.SendJSONResponse(w, jRet)
	return
}

// PutNotificationSettingsHandler puts notification settings.
func PutNotificationSettingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	q := r.URL.Query()

	notificationTypeID := ps.ByName("notificationTypeID")
	if len(notificationTypeID) < 1 {
		util.SendBadRequest(w, "missing notificationTypeID as url param")
		return
	}

	ntID, err := strconv.Atoi(notificationTypeID)
	if err != nil {
		util.SendBadRequest(w, "invalid notification_type_id as query param")
		return
	}

	medium := notifications.Medium(q.Get("medium"))
	if !medium.IsValid() {
		util.SendBadRequest(w, "invalid medium as query param")
		return
	}

	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeNotifications}, &oauth2.AuthorizerAPICall{
		Method:    "PUT",
		RoutePath: "notifications/types/:notificationTypeID",
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	err = notifications.InsertNotificationSettings(db, &notifications.InsertNotificationSettingsRequest{
		UserID: *userID,
		TypeID: uint64(ntID),
		Medium: medium,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == util.PgErrorForeignKeyViolation {
				util.SendNotFoundWithReason(w, "unknown notificationTypeID as url param")
			}
			return
		}
		handleISE(w, fmt.Errorf("error inserting notification settings"))
		return
	}

	util.SendNoContent(w)
	return
}

// DeleteNotificationSettingHandler deletes a notification for a user.
func DeleteNotificationSettingHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	medium := notifications.Medium(r.URL.Query().Get("medium"))
	if !medium.IsValid() {
		util.SendBadRequest(w, "invalid medium as query param")
		return
	}

	notificationTypeID := ps.ByName("notificationTypeID")
	if len(notificationTypeID) < 1 {
		util.SendBadRequest(w, "missing notificationTypeID as url param")
		return
	}

	ntID, err := strconv.Atoi(notificationTypeID)
	if err != nil {
		util.SendBadRequest(w, "invalid notification_type_id as query param")
		return
	}

	userID, rdP, sess := authorizer(w, r, []oauth2.Scope{oauth2.ScopeNotifications}, &oauth2.AuthorizerAPICall{
		Method:    "DELETE",
		RoutePath: "notifications/types/:notificationTypeID",
	})
	if (userID == nil || rdP == nil) && sess == nil {
		return
	}

	err = notifications.DeleteNotificationSetting(db, &notifications.DeleteNotificationSettingRequest{
		UserID: *userID,
		TypeID: uint64(ntID),
		Medium: medium,
	})
	if err != nil {
		handleISE(w, fmt.Errorf("error deleting a notification setting: %w", err))
		return
	}

	util.SendNoContent(w)
	return
}
