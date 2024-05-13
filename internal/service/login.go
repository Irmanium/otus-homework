package service

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

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
