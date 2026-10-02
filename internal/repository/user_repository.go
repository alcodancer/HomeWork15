// internal/repository/user_repository.go
package repository

import (
	"context"

	"example/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]*domain.User, error)
}

type OrderRepository interface {
	Create(ctx context.Context, o *domain.Order) error
	ListByUser(ctx context.Context, userID string) ([]*domain.Order, error)
}

// TxRunner — нужен для транзакций в сервисе
type TxRunner interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context, ur UserRepository, or OrderRepository) error) error
}