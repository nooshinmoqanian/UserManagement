package user

import (
	"context"
	"strings"

	"github.com/you/otp-auth/internal/domain"
	"gorm.io/gorm"
)

type GormRepo struct{ db *gorm.DB }

func NewGormRepo(db *gorm.DB) *GormRepo { return &GormRepo{db: db} }

func (r *GormRepo) UpsertByPhone(ctx context.Context, phone string) (domain.User, error) {
	phone = strings.TrimSpace(phone)
	u := domain.User{Phone: phone}
	if err := r.db.WithContext(ctx).FirstOrCreate(&u, "phone = ?", phone).Error; err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (r *GormRepo) Get(ctx context.Context, phone string) (domain.User, bool, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).First(&u, "phone = ?", phone).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return domain.User{}, false, nil }
		return domain.User{}, false, err
	}
	return u, true, nil
}

func (r *GormRepo) List(ctx context.Context, q domain.ListUsersQuery) (domain.PaginatedUsers, error) {
	if q.Limit <= 0 { q.Limit = 10 }
	if q.Page  <= 0 { q.Page  = 1 }
	offset := (q.Page - 1) * q.Limit

	tx := r.db.WithContext(ctx).Model(&domain.User{})
	if strings.TrimSpace(q.Search) != "" {
		tx = tx.Where("phone ILIKE ?", "%"+q.Search+"%")
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil { return domain.PaginatedUsers{}, err }

	var items []domain.User
	if err := tx.Order("registration DESC").Limit(q.Limit).Offset(offset).Find(&items).Error; err != nil {
		return domain.PaginatedUsers{}, err
	}

	return domain.PaginatedUsers{ Items: items, Page: q.Page, Limit: q.Limit, TotalItems: int(total) }, nil
}
