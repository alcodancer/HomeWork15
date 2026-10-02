// internal/storage/postgres/numeric.go
package postgres

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func numericFromFloat(f float64) pgtype.Numeric {
	var n pgtype.Numeric
	_ = n.Scan(f)
	return n
}