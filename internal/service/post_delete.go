package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

func (s *Service) rebuildFriendsFeedCache(userID string) {
	go func(ctx context.Context, userID string) {
		friends, err := s.userRepo.GetFriends(ctx, userID)
		if err != nil {
			s.Logger.Warn(err)
			return
		}

		for _, friend := range friends {
			s.rebuildUserFeedCache(friend)
		}
	}(context.Background(), userID)
}

func (s *Service) deletePost(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)

	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid post id"})
	}

	postUserID, _, _, err := s.userRepo.GetPost(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "post not found"})
		}

		return err
	}
	if userID != postUserID {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "you can edit only your posts"})
	}

	if err = s.userRepo.DeletePost(c.Request().Context(), id); err != nil {
		return err
	}

	s.rebuildFriendsFeedCache(userID)

	return c.String(http.StatusOK, "")
}
