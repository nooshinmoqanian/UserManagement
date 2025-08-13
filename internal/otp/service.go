package otp

import (
	"context"
	"time"
)

type Service interface {
	Request(ctx context.Context, phone string) (code string, expiresAt time.Time, err error)
	Verify(ctx context.Context, phone, code string) (bool, error)
}
