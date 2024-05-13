package service

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

type Message struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

type GetDialogResp []Message

func (s *Service) getDialog(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)
	interlocutorID := c.Param("user_id")

	_, err := s.userRepo.GetUser(c.Request().Context(), interlocutorID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "user not found"})
		}

		return err
	}

	messages, err := s.dialogRepo.GetDialog(c.Request().Context(), userID, interlocutorID)
	if err != nil {
		return err
	}

	resp := make(GetDialogResp, 0, len(messages))
	for _, message := range messages {
		resp = append(resp, Message{
			From: message.From,
			To:   message.To,
			Text: message.Text,
		})
	}

	return c.JSON(http.StatusOK, resp)
}
