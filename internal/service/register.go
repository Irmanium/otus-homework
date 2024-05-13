package service

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

type RegisterReq struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
	Password   string `json:"password"`
}

type RegisterResp struct {
	UserID string `json:"user_id"`
}

func (s *Service) register(c echo.Context) error {
	req := new(RegisterReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid request"})
	}

	birthdate, err := time.Parse(time.DateOnly, req.Birthdate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid birthdate"})
	}

	passwordHash, err := generatePasswordHash(req.Password)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	err = s.userRepo.Register(c.Request().Context(), domain.FullUser{
		UserProfile: domain.UserProfile{
			ID:         id,
			FirstName:  req.FirstName,
			SecondName: req.SecondName,
			Birthdate:  birthdate,
			Biography:  req.Biography,
			City:       req.City,
		},
		PasswordHash: passwordHash,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, RegisterResp{UserID: id})
}
