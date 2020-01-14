package gift_cards

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/iamtheyammer/canvascbl/backend/src/db/services"
	"github.com/iamtheyammer/canvascbl/backend/src/util"
	"github.com/pkg/errors"
	"time"
)

type ListRequest struct {
	IDs        []uint64
	ClaimCodes []string
	RedeemedBy uint64

	ValidOnly bool
	Limit     uint64
	Offset    uint64
}

type GiftCard struct {
	ID         uint64
	ClaimCode  string
	ValidFor   uint64
	ExpiresAt  time.Time
	RedeemedAt time.Time
	RedeemedBy uint64
	InsertedAt time.Time
}

func List(db services.DB, req *ListRequest) (*[]GiftCard, error) {
	q := util.Sq.
		Select(
			"id",
			"claim_code",
			"valid_for",
			"expires_at",
			"redeemed_at",
			"redeemed_by",
			"inserted_at",
		).
		From("gift_cards")

	if len(req.IDs) > 0 {
		q = q.Where(sq.Eq{"id": req.IDs})
	}

	if len(req.ClaimCodes) > 0 {
		q = q.Where(sq.Eq{"claim_code": req.ClaimCodes})
	}

	if req.RedeemedBy != 0 {
		q = q.Where(sq.Eq{"redeemed_by": req.RedeemedBy})
	}

	if req.ValidOnly {
		q = q.
			Where(sq.Eq{"redeemed_at": nil}).
			Where("expires_at IS NULL OR expires_at < NOW()")
	}

	if req.Limit != 0 {
		q = q.Limit(req.Limit)
	}

	if req.Offset != 0 {
		q = q.Offset(req.Offset)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building list gift cards sql")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error executing list gift cards sql")
	}

	defer rows.Close()

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
			return nil, errors.Wrap(err, "error scanning list gift cards sql")
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
