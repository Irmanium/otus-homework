package service

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

type GetPostResp struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	UserID string `json:"author_user_id"`
}

func (s *Service) getPost(c echo.Context) error {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid post id"})
	}

	userID, text, _, err := s.userRepo.GetPost(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "post not found"})
		}

		return err
	}

	return c.JSON(http.StatusOK, GetPostResp{ID: id, Text: text, UserID: userID})
}
