package gift_cards

import (
	"database/sql"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type InsertRequest struct {
	ClaimCodes []string
	ValidFor   uint64
	ExpiresAt  *time.Time
}

func Insert(db services.DB, req *InsertRequest) (*[]GiftCard, error) {
	q := util.Sq.
		Insert("gift_cards").
		Columns("claim_code", "valid_for", "expires_at").
		Suffix("RETURNING id, claim_code, valid_for, expires_at, redeemed_at, redeemed_by, inserted_at")

	for _, cc := range req.ClaimCodes {
		q = q.Values(cc, req.ValidFor, req.ExpiresAt)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building insert gift cards sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing insert gift cards sql")
	}

	var gcs []GiftCard
	for rows.Next() {
		var (
			gc                    GiftCard
			expiresAt, redeemedAt sql.NullTime
			redeemedBy            sql.NullInt64
		)
		err := rows.Scan(
			&gc.ID,
			&gc.ClaimCode,
			&gc.ValidFor,
			&expiresAt,
			&redeemedAt,
			&redeemedBy,
			&gc.InsertedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning gift cards")
		}

		if expiresAt.Valid {
			gc.ExpiresAt = expiresAt.Time
		}

		if redeemedAt.Valid {
			gc.RedeemedAt = redeemedAt.Time
		}

		if redeemedBy.Valid {
			gc.RedeemedBy = uint64(redeemedBy.Int64)
		}

		gcs = append(gcs, gc)
	}

	return &gcs, nil
}
