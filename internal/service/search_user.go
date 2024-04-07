package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SearchUserResp []GetUserResp

func (s *Service) searchUser(c echo.Context) error {
	firstName := c.QueryParam("first_name")
	secondName := c.QueryParam("second_name")

	users, err := s.r.SearchUser(c.Request().Context(), firstName, secondName)
	if err != nil {
		return err
	}

	resp := make(SearchUserResp, 0, len(users))
	for _, user := range users {
		resp = append(resp, convertDomainProfileToResp(user))
	}

	return c.JSON(http.StatusOK, resp)
}
