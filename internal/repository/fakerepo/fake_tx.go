// internal/repository/fakerepo/fake_tx.go
package fakerepo

import (
	"context"

	"example/internal/repository"
)

// FakeTxRunner вызывает fn с теми же репозиториями, без реальной транзакции.
type FakeTxRunner struct {
	Users  repository.UserRepository
	Orders repository.OrderRepository
}

func (t *FakeTxRunner) RunInTx(
	ctx context.Context,
	fn func(ctx context.Context, ur repository.UserRepository, or repository.OrderRepository) error,
) error {
	return fn(ctx, t.Users, t.Orders)
}