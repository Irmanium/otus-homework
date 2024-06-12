package service

import (
	"errors"
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

type LoginReq struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token string `json:"token"`
}

func (s *Service) login(c echo.Context) error {
	req := new(LoginReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid request"})
	}

	_, err := uuid.Parse(req.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid id"})
	}

	passwordHash, err := s.userRepo.GetPassword(c.Request().Context(), req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "user not found"})
		}

		return err
	}

	err = checkPassword(passwordHash, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrResp{Message: "incorrect password"})
	}

	token, err := s.generateToken(req.ID)
	if err != nil {
		return nil
	}

	return c.JSON(http.StatusOK, LoginResp{Token: token})
}
