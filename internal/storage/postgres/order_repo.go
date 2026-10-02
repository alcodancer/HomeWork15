// internal/storage/postgres/order_repo.go
package postgres

import (
	"context"

	"github.com/google/uuid"

	"example/internal/domain"
	"example/internal/storage/postgres/sqlc"
)

type OrderRepo struct {
	q *sqlc.Queries
}

func NewOrderRepo(q *sqlc.Queries) *OrderRepo {
	return &OrderRepo{q: q}
}

func (r *OrderRepo) Create(ctx context.Context, o *domain.Order) error {
	amount, _ := o.Amount, error(nil)
	row, err := r.q.CreateOrder(ctx, sqlc.CreateOrderParams{
		UserID: o.UserID,
		Amount: numericFromFloat(amount),
		Note:   o.Note,
	})
	if err != nil {
		return err
	}
	*o = *mapOrder(row)
	return nil
}

func (r *OrderRepo) ListByUser(ctx context.Context, userID string) ([]*domain.Order, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	rows, err := r.q.ListOrdersByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.Order, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapOrder(r))
	}
	return out, nil
}

func mapOrder(row sqlc.Order) *domain.Order {
	f, _ := row.Amount.Float64Value()
	return &domain.Order{
		ID:        row.ID,
		UserID:    row.UserID,
		Amount:    f.Float64,
		Status:    row.Status,
		Note:      row.Note,
		CreatedAt: row.CreatedAt.Time,
	}
}