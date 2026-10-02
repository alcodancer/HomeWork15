// internal/service/user_service.go (фрагмент)
package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"example/internal/domain"
	"example/internal/repository"
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

type UserService struct {
	repo    repository.UserRepository
	orders  repository.OrderRepository
	history repository.HistoryRepository
	tx      repository.TxRunner
}

func NewUserService(
	repo repository.UserRepository,
	orders repository.OrderRepository,
	history repository.HistoryRepository,
	tx repository.TxRunner,
) *UserService {
	return &UserService{repo: repo, orders: orders, history: history, tx: tx}
}

// CreateUserWithOrder — атомарно создаёт пользователя и его первый заказ.
func (s *UserService) CreateUserWithOrder(ctx context.Context, name, email string, amount float64, note *string) (*domain.User, *domain.Order, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))

	if name == "" {
		return nil, nil, domain.ErrInvalidName
	}
	if !emailRe.MatchString(email) {
		return nil, nil, domain.ErrInvalidEmail
	}
	if amount <= 0 {
		return nil, nil, domain.ErrInvalidOrder
	}

	var (
		user  domain.User
		order domain.Order
	)

	err := s.tx.RunInTx(ctx, func(ctx context.Context, ur repository.UserRepository, or repository.OrderRepository) error {
		// 1) создаём пользователя
		user = domain.User{Name: name, Email: email}
		if err := ur.Create(ctx, &user); err != nil {
			return err
		}

		// 2) создаём заказ
		order = domain.Order{UserID: user.ID, Amount: amount, Note: note}
		return or.Create(ctx, &order)
	})
	if err != nil {
		return nil, nil, err
	}

	s.record(ctx, "create_with_order", user.ID.String(), map[string]any{
		"user":  user,
		"order": order,
	})

	return &user, &order, nil
}

func (s *UserService) record(ctx context.Context, action, id string, data any) {
	if s.history == nil {
		return
	}
	_ = s.history.Record(ctx, &repository.HistoryEntry{
		Action: action,
		UserID: id,
		Data:   data,
	})
}