package service

import (
	"context"
	"time"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"otus-homework/internal/domain"
)

type UserRepo interface {
	Register(ctx context.Context, user domain.FullUser) error
	GetUser(ctx context.Context, id string) (domain.UserProfile, error)
	GetPassword(ctx context.Context, id string) (string, error)
	SearchUser(ctx context.Context, firstName, secondName string) ([]domain.UserProfile, error)
}

type DialogRepo interface {
	SendMessage(ctx context.Context, id string, message domain.Message) error
	GetDialog(ctx context.Context, userID, interlocutorID string) ([]domain.Message, error)
}

type Service struct {
	*echo.Echo

	userRepo      UserRepo
	dialogRepo    DialogRepo
	port          string
	jwtSecret     string
	tokenTTLHours int
}

type ErrResp struct {
	Message string `json:"message"`
}

func New(userRepo UserRepo, dialogRepo DialogRepo, port, jwtSecret string, tokenTTLHours int) *Service {
	e := echo.New()

	return &Service{
		Echo:          e,
		userRepo:      userRepo,
		dialogRepo:    dialogRepo,
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
	r.GET("/user/search", s.searchUser)
	r.POST("/dialog/:user_id/send", s.sendMessage)
	r.GET("/dialog/:user_id/list", s.getDialog)

	s.Logger.Fatal(s.Start(":" + s.port))
}

func convertDomainProfileToResp(user domain.UserProfile) GetUserResp {
	return GetUserResp{
		ID:         user.ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Birthdate:  user.Birthdate.Format(time.DateOnly),
		Biography:  user.Biography,
		City:       user.City,
	}
}
