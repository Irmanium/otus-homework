package service

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

func (s *Service) postFeed(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)
	friendIDs, err := s.userRepo.GetFriends(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	cancel := make(chan struct{})
	feed, err := s.liveFeedRepo.GetFeed(friendIDs, cancel)
	if err != nil {
		return err
	}

	go func() {
		for {
			mt, _, err := ws.ReadMessage()
			if err != nil || mt == websocket.CloseMessage {
				cancel <- struct{}{}
			}
		}
	}()

	for post := range feed {
		err := ws.WriteMessage(websocket.TextMessage, post)
		if err != nil {
			return err
		}
	}

	return nil
}
