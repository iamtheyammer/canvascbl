package gift_cards

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
)

type UpdateRequest struct {
	Where ListRequest
	Set   GiftCard
}

func Update(db services.DB, req *UpdateRequest) error {
	q := util.Sq.
		Update("gift_cards")

	if len(req.Where.IDs) > 0 {
		q = q.Where(sq.Eq{"id": req.Where.IDs})
	}

	if len(req.Where.ClaimCodes) > 0 {
		q = q.Where(sq.Eq{"claim_code": req.Where.ClaimCodes})
	}

	if req.Where.RedeemedBy != 0 {
		q = q.Where(sq.Eq{"redeemed_by": req.Where.RedeemedBy})
	}

	if req.Where.ValidOnly {
		q = q.
			Where("expires_at < NOW()").
			Where(sq.Eq{"redeemed_at": nil})
	}

	if len(req.Set.ClaimCode) > 1 {
		q = q.Set("claim_code", req.Set.ClaimCode)
	}

	if req.Set.ValidFor != 0 {
		q = q.Set("valid_for", req.Set.ValidFor)
	}

	if !req.Set.ExpiresAt.IsZero() {
		q = q.Set("expires_at", req.Set.ExpiresAt)
	}

	if !req.Set.RedeemedAt.IsZero() {
		q = q.Set("redeemed_at", req.Set.RedeemedAt)
	}

	if req.Set.RedeemedBy != 0 {
		q = q.Set("redeemed_by", req.Set.RedeemedBy)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return errors.Wrap(err, "error building update gift cards sql")
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return errors.Wrap(err, "error executing update gift cards sql")
	}

	return nil
}
