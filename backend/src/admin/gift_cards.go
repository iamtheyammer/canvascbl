package admin

import (
	"encoding/json"
	"github.com/iamtheyammer/canvascbl/backend/src/db"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services/gift_cards"
	"github.com/iamtheyammer/canvascbl/backend/src/middlewares"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type generateGiftCardsResponse struct {
	GiftCards []giftCard `json:"gift_cards"`
}

type giftCard struct {
	ID         uint64     `json:"id"`
	ClaimCode  string     `json:"claimCode"`
	ValidFor   uint64     `json:"validFor"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	RedeemedAt *time.Time `json:"redeemedAt,omitempty"`
	RedeemedBy *uint64    `json:"redeemedBy,omitempty"`
	InsertedAt time.Time  `json:"insertedAt"`
}

func GenerateGiftCardsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	quantity := r.URL.Query().Get("quantity")
	if len(quantity) < 1 {
		util.SendBadRequest(w, "missing quantity as query param")
		return
	}

	qty, err := strconv.Atoi(quantity)
	if err != nil {
		util.SendBadRequest(w, "error converting quantity as query param to an integer")
		return
	}

	validFor := r.URL.Query().Get("valid_for")
	if len(validFor) < 1 {
		util.SendBadRequest(w, "missing valid_for as query param")
		return
	}

	vf, err := strconv.Atoi(validFor)
	if err != nil {
		util.SendBadRequest(w, "error converting valid_for as query param to an integer")
		return
	}

	session := middlewares.Session(w, r, true)
	if session == nil {
		return
	}

	if middlewares.IsAdmin(w, r, session) {
		return
	}

	claimCodes := util.GenerateGiftCardClaimCodes(qty)

	gcs, err := db.InsertGiftCards(&gift_cards.InsertRequest{
		ClaimCodes: claimCodes,
		ValidFor:   uint64(vf),
		ExpiresAt:  nil,
	})
	if err != nil {
		util.HandleError(errors.Wrap(err, "error inserting gift cards"))
		util.SendInternalServerError(w)
		return
	}

	var giftCards []giftCard
	for _, gc := range *gcs {
		giftCrd := giftCard{
			ID:         gc.ID,
			ClaimCode:  gc.ClaimCode,
			ValidFor:   gc.ValidFor,
			InsertedAt: gc.InsertedAt,
		}

		if !gc.ExpiresAt.IsZero() {
			giftCrd.ExpiresAt = &gc.ExpiresAt
		}

		if !gc.RedeemedAt.IsZero() {
			giftCrd.RedeemedAt = &gc.RedeemedAt
		}

		if gc.RedeemedBy != 0 {
			giftCrd.RedeemedBy = &gc.RedeemedBy
		}

		giftCards = append(giftCards, giftCrd)
	}

	jGiftCards, err := json.Marshal(&giftCards)
	if err != nil {
		util.HandleError(errors.Wrap(err, "error marshaling generate gift card response"))
		util.SendInternalServerError(w)
		return
	}

	util.SendJSONResponse(w, jGiftCards)
	return
}
