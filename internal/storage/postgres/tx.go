// internal/storage/postgres/tx.go  (дополнение)
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"example/internal/repository"
	"example/internal/storage/postgres/sqlc"
)

type TxRunner struct {
	pool *pgxpool.Pool
}

func NewTxRunner(pool *pgxpool.Pool) *TxRunner {
	return &TxRunner{pool: pool}
}

func (t *TxRunner) RunInTx(
	ctx context.Context,
	fn func(ctx context.Context, ur repository.UserRepository, or repository.OrderRepository) error,
) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	ur := NewUserRepo(q)
	or := NewOrderRepo(q)

	if err := fn(ctx, ur, or); err != nil {
		return err
	}
	return tx.Commit(ctx)
}