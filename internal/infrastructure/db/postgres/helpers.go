package postgres

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func timestamptzFromTime(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	_ = ts.Scan(t)
	return ts
}

func timeFromTimestamptz(ts pgtype.Timestamptz) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}
