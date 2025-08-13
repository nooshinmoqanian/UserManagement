package httpserver

import (
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/you/otp-auth/docs" // swagger docs

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/you/otp-auth/internal/auth"
	"github.com/you/otp-auth/internal/domain"
	"github.com/you/otp-auth/internal/otp"
	"github.com/you/otp-auth/internal/user"
)

type Config struct {
	JWTSecret     []byte
	UserRepo      user.Repository
	OTPSvc        otp.Service
	EnableSwagger bool
}

type Server struct {
	engine *gin.Engine
	h      *Handler
}

type Handler struct {
	user user.Repository
	otp  otp.Service
	jwt  *auth.JWT
}

func New(cfg Config) *Server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), requestLogger())

	h := &Handler{
		user: cfg.UserRepo,
		otp:  cfg.OTPSvc,
		jwt:  auth.New(cfg.JWTSecret),
	}

	// Health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	// Swagger UI
	if cfg.EnableSwagger {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Auth
	authg := r.Group("/auth")
	{
		authg.POST("/request-otp", h.RequestOTP)
		authg.POST("/verify-otp", h.VerifyOTP)
	}

	// Protected users routes
	users := r.Group("/users", h.authMW)
	{
		users.GET("/me", h.GetMe)
		users.GET("/:phone", h.GetByPhone)
		users.GET("", h.ListUsers)
	}

	return &Server{engine: r, h: h}
}

func (s *Server) Run(addr string) error { return s.engine.Run(addr) }

// ===== Handlers + Swagger =====

// @Summary Request OTP
// @Description Generate OTP for a phone number. Rate limited to 3 per 10 minutes. OTP expires in 2 minutes.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.OTPRequest true "Phone number"
// @Success 200 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /auth/request-otp [post]
func (h *Handler) RequestOTP(c *gin.Context) {
	var req domain.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Phone) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone is required"})
		return
	}
	code, exp, err := h.otp.Request(c.Request.Context(), req.Phone)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[OTP] phone=%s code=%s expires=%s", req.Phone, code, exp.Format(time.RFC3339))
	c.JSON(http.StatusOK, gin.H{"message": "otp generated", "expiresAt": exp.Format(time.RFC3339)})
}

// @Summary Verify OTP (login/register)
// @Description Verify OTP; creates user if not exists, returns JWT.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.OTPVerify true "phone+otp"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /auth/verify-otp [post]
func (h *Handler) VerifyOTP(c *gin.Context) {
	var req domain.OTPVerify
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "phone and otp are required"})
		return
	}
	ok, err := h.otp.Verify(c.Request.Context(), req.Phone, req.OTP)
	if err != nil || !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired otp"})
		return
	}
	u, err := h.user.UpsertByPhone(c.Request.Context(), req.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	token, err := h.jwt.Sign(u.Phone, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot issue token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary Get current user
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.User
// @Router /users/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	phone := c.GetString("phone")
	u, ok, _ := h.user.Get(c.Request.Context(), phone)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// @Summary Get user by phone
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param phone path string true "phone"
// @Success 200 {object} domain.User
// @Failure 404 {object} map[string]string
// @Router /users/{phone} [get]
func (h *Handler) GetByPhone(c *gin.Context) {
	phone := c.Param("phone")
	u, ok, _ := h.user.Get(c.Request.Context(), phone)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// @Summary List users (pagination + search)
// @Tags users
// @Produce json
// @Security BearerAuth
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param search query string false "search by phone"
// @Success 200 {object} domain.PaginatedUsers
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	var q domain.ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad query"})
		return
	}
	res, err := h.user.List(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, res)
}
