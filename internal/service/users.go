package service

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

type GetUserResp struct {
	ID         string `json:"id"`
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	Birthdate  string `json:"birthdate"`
	Biography  string `json:"biography"`
	City       string `json:"city"`
}

func (s *Service) getUser(c echo.Context) error {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid id"})
	}

	user, err := s.userRepo.GetUser(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "user not found"})
		}

		return err
	}

	return c.JSON(http.StatusOK, convertDomainProfileToResp(user))
}

type SearchUserResp []GetUserResp

func (s *Service) searchUser(c echo.Context) error {
	firstName := c.QueryParam("first_name")
	secondName := c.QueryParam("second_name")

	users, err := s.userRepo.SearchUser(c.Request().Context(), firstName, secondName)
	if err != nil {
		return err
	}

	resp := make(SearchUserResp, 0, len(users))
	for _, user := range users {
		resp = append(resp, convertDomainProfileToResp(user))
	}

	return c.JSON(http.StatusOK, resp)
}
