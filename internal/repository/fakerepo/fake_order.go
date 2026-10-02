// internal/repository/fakerepo/fake_order.go
package fakerepo

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"example/internal/domain"
)

type FakeOrderRepo struct {
	mu     sync.Mutex
	Orders map[string][]*domain.Order
	Err    error
}

func NewFakeOrderRepo() *FakeOrderRepo {
	return &FakeOrderRepo{Orders: map[string][]*domain.Order{}}
}

func (r *FakeOrderRepo) Create(_ context.Context, o *domain.Order) error {
	if r.Err != nil {
		return r.Err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	o.ID = uuid.New()
	cp := *o
	r.Orders[o.UserID.String()] = append(r.Orders[o.UserID.String()], &cp)
	return nil
}

func (r *FakeOrderRepo) ListByUser(_ context.Context, userID string) ([]*domain.Order, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return r.Orders[userID], nil
}