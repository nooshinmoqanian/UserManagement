package user

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/you/otp-auth/internal/domain"
)

type MemoryRepo struct {
	mu    sync.RWMutex
	users map[string]domain.User // key: phone
}

func NewMemoryRepo() *MemoryRepo { return &MemoryRepo{users: make(map[string]domain.User)} }

func (s *MemoryRepo) UpsertByPhone(ctx context.Context, phone string) (domain.User, error) {
	s.mu.Lock(); defer s.mu.Unlock()
	if u, ok := s.users[phone]; ok {
		return u, nil
	}
	u := domain.User{ Phone: phone, Registration: time.Now() }
	s.users[phone] = u
	return u, nil
}

func (s *MemoryRepo) Get(ctx context.Context, phone string) (domain.User, bool, error) {
	s.mu.RLock(); defer s.mu.RUnlock()
	u, ok := s.users[phone]
	return u, ok, nil
}

func (s *MemoryRepo) List(ctx context.Context, q domain.ListUsersQuery) (domain.PaginatedUsers, error) {
	s.mu.RLock()
	all := make([]domain.User, 0, len(s.users))
	for _, u := range s.users { all = append(all, u) }
	s.mu.RUnlock()

	if strings.TrimSpace(q.Search) != "" {
		needle := strings.ToLower(q.Search)
		filtered := all[:0]
		for _, u := range all {
			if strings.Contains(strings.ToLower(u.Phone), needle) {
				filtered = append(filtered, u)
			}
		}
		all = filtered
	}

	sort.Slice(all, func(i, j int) bool { return all[i].Registration.After(all[j].Registration) })

	if q.Limit <= 0 { q.Limit = 10 }
	if q.Page <= 0 { q.Page = 1 }
	start := (q.Page - 1) * q.Limit
	if start > len(all) { start = len(all) }
	end := start + q.Limit
	if end > len(all) { end = len(all) }

	return domain.PaginatedUsers{
		Items: all[start:end],
		Page: q.Page,
		Limit: q.Limit,
		TotalItems: len(all),
	}, nil
}
