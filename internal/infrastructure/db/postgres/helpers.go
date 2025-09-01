package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func timestamptzFromTime(t time.Time) pgtype.Timestamptz {
	var ts pgtype.Timestamptz
	ts.Scan(t)
	return ts
}

func timeFromTimestamptz(ts pgtype.Timestamptz) time.Time {
	if ts.Valid {
		return ts.Time
	}
	return time.Time{}
}

func numericFromFloat64(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	n.Scan(f)
	return n
}

func float64FromNumeric(n pgtype.Numeric) float64 {
	f, _ := n.Float64Value()
	return f.Float64
}
