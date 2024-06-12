package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

type UpdatePostReq struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (s *Service) editPostInUserFeedCache(userID, postID, text string) {
	go func(ctx context.Context, userID, postID, text string) {
		posts, fromCache, err := s.getPostsCached(ctx, userID)
		if err != nil {
			s.Logger.Warn(err)
			return
		}

		if !fromCache {
			return
		}

		for i, post := range posts {
			if post.ID == postID {
				posts[i].Text = text
				break
			}
		}

		err = s.feedCache.PutFeedToCache(ctx, userID, posts)
		if err != nil {
			s.Logger.Warn(err)
		}
	}(context.Background(), userID, postID, text)
}

func (s *Service) editPostInFriendsFeedCache(userID, postID, text string) {
	go func(ctx context.Context, userID string) {
		friends, err := s.userRepo.GetFriends(ctx, userID)
		if err != nil {
			s.Logger.Warn(err)
			return
		}

		for _, friend := range friends {
			s.editPostInUserFeedCache(friend, postID, text)
		}
	}(context.Background(), userID)
}

func (s *Service) updatePost(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)

	req := new(UpdatePostReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid request"})
	}

	_, err := uuid.Parse(req.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid post id"})
	}

	postUserID, _, _, err := s.userRepo.GetPost(c.Request().Context(), req.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return c.JSON(http.StatusNotFound, ErrResp{Message: "post not found"})
		}

		return err
	}
	if userID != postUserID {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "you can edit only your posts"})
	}

	if err = s.userRepo.UpdatePost(c.Request().Context(), req.ID, req.Text); err != nil {
		return err
	}

	s.editPostInFriendsFeedCache(userID, req.ID, req.Text)

	return c.String(http.StatusOK, "")
}
