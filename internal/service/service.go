package service

import (
	"context"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"otus-homework/internal/domain"
)

type Repo interface {
	Register(ctx context.Context, user domain.FullUser) error
	GetUser(ctx context.Context, id string) (domain.UserProfile, error)
	GetPassword(ctx context.Context, id string) (string, error)
}

type Service struct {
	*echo.Echo

	r             Repo
	port          string
	jwtSecret     string
	tokenTTLHours int
}

type ErrResp struct {
	Message string `json:"message"`
}

func New(repo Repo, port, jwtSecret string, tokenTTLHours int) *Service {
	e := echo.New()

	return &Service{
		Echo:          e,
		r:             repo,
		port:          port,
		jwtSecret:     jwtSecret,
		tokenTTLHours: tokenTTLHours,
	}
}

func (s *Service) StartService() {
	s.Use(middleware.Logger())
	s.Use(middleware.Recover())

	s.POST("/login", s.login)
	s.POST("/user/register", s.register)

	r := s.Group("")
	r.Use(echojwt.WithConfig(s.getJWTConfig()))
	r.GET("/user/get/:id", s.getUser)

	s.Logger.Fatal(s.Start(":" + s.port))
}
