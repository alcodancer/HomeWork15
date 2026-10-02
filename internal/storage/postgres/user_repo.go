// internal/storage/postgres/user_repo.go
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"example/internal/domain"
	"example/internal/storage/postgres/sqlc"
)

type UserRepo struct {
	q *sqlc.Queries
}

func NewUserRepo(q *sqlc.Queries) *UserRepo {
	return &UserRepo{q: q}
}

func mapUser(row sqlc.User) *domain.User {
	return &domain.User{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		Phone:     row.Phone,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	row, err := r.q.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrUserExists
		}
		return err
	}
	*u = *mapUser(row)
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	row, err := r.q.GetUserByID(ctx, uid)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapUser(row), nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return mapUser(row), nil
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	row, err := r.q.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrUserNotFound
	}
	if err != nil {
		if isUniqueViolation(err) {
			return domain.ErrUserExists
		}
		return err
	}
	*u = *mapUser(row)
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrUserNotFound
	}
	n, err := r.q.DeleteUser(ctx, uid)
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepo) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.q.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.User, 0, len(rows))
	for _, r := range rows {
		out = append(out, mapUser(r))
	}
	return out, nil
}

// isUniqueViolation определяет ошибку уникальности postgres (код 23505)
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}