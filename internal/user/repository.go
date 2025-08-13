package user

import (
	"context"
	"github.com/you/otp-auth/internal/domain"
)

type Repository interface {
	UpsertByPhone(ctx context.Context, phone string) (domain.User, error)
	Get(ctx context.Context, phone string) (domain.User, bool, error)
	List(ctx context.Context, q domain.ListUsersQuery) (domain.PaginatedUsers, error)
}
