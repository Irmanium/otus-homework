package service

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

type SendMessageReq struct {
	Text string `json:"text"`
}

func (s *Service) sendMessage(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)
	recipientID := c.Param("user_id")

	req := new(SendMessageReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid request"})
	}

	if _, err := s.userRepo.GetUser(c.Request().Context(), recipientID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "user not found"})
		}

		return err
	}

	if err := s.dialogRepo.SendMessage(c.Request().Context(), uuid.New().String(), domain.Message{
		From: userID,
		To:   recipientID,
		Text: req.Text,
	}); err != nil {
		return err
	}

	return c.String(http.StatusOK, "")
}
