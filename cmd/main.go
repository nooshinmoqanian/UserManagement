// @title           OTP Auth Service API
// @version         1.0
// @description     OTP-based login/registration with rate limiting and user management.
// @schemes         http
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/you/otp-auth/internal/cache"
	"github.com/you/otp-auth/internal/config"
	"github.com/you/otp-auth/internal/db"
	"github.com/you/otp-auth/internal/domain"
	"github.com/you/otp-auth/internal/otp"
	"github.com/you/otp-auth/internal/transport/httpserver"
	"github.com/you/otp-auth/internal/user"
)

func main() {
	cfg := config.Load()

	gdb := db.ConnectPostgres()
	if err := gdb.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("AutoMigrate error: %v", err)
	}

	rdb := cache.ConnectRedis()

	userRepo := user.NewGormRepo(gdb)
	otpSvc := otp.NewRedis(otp.RedisConfig{
		RDB:        rdb,
		ExpiresIn:  cfg.OTPExpire,
		MaxRequests: cfg.OTPMaxReq,
		Window:     cfg.OTPWindow,
	})

	s := httpserver.New(httpserver.Config{
		JWTSecret:     []byte(cfg.JWTSecret),
		UserRepo:      userRepo,
		OTPSvc:        otpSvc,
		EnableSwagger: true,
	})

	log.Printf("Listening on :%s", cfg.Port)
	log.Fatal(s.Run(":" + cfg.Port))
}
