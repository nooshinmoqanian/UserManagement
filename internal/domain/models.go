package domain

import "time"

type User struct {
	Phone        string    `json:"phone" gorm:"primaryKey;type:text"`
	Registration time.Time `json:"registrationDate" gorm:"not null;default:now()"`
}

type OTPRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type OTPVerify struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp"   binding:"required"`
}

type ListUsersQuery struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=10"`
	Search string `form:"search"`
}

type PaginatedUsers struct {
	Items      []User `json:"items"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalItems int    `json:"totalItems"`
}
