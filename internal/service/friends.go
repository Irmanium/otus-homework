package service

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

func (s *Service) setFriend(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)
	friendID := c.Param("user_id")
	if userID == friendID {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "you can`t add yourself to friends"})
	}

	_, err := uuid.Parse(friendID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid friend id"})
	}

	if _, err := s.userRepo.GetUser(c.Request().Context(), friendID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "user not found"})
		}

		return err
	}

	if err := s.userRepo.SetFriend(c.Request().Context(), [2]string{userID, friendID}); err != nil {
		return err
	}

	s.rebuildUserFeedCache(userID)
	s.rebuildUserFeedCache(friendID)

	return c.String(http.StatusOK, "")
}

func (s *Service) deleteFriend(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)
	friendID := c.Param("user_id")
	if userID == friendID {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "you can`t delete yourself from friends"})
	}

	_, err := uuid.Parse(friendID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid friend id"})
	}

	if _, err := s.userRepo.GetUser(c.Request().Context(), friendID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "user not found"})
		}

		return err
	}

	if err := s.userRepo.DeleteFriend(c.Request().Context(), [2]string{userID, friendID}); err != nil {
		return err
	}

	s.rebuildUserFeedCache(userID)
	s.rebuildUserFeedCache(friendID)

	return c.String(http.StatusOK, "")
}
