package service

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

func (s *Service) getPostsCached(ctx context.Context, userID string) (posts []domain.Post, fromCache bool, err error) {
	posts, err = s.feedCache.GetFeedFromCache(ctx, userID)
	if err == nil {
		return posts, true, err
	} else {
		s.Logger.Warn(err)
	}

	posts, err = s.userRepo.GetFeed(ctx, userID)
	if err != nil {
		return nil, false, err
	}

	go func(ctx context.Context, userID string, posts []domain.Post) {
		err = s.feedCache.PutFeedToCache(ctx, userID, posts)
		if err != nil {
			s.Logger.Warn(err)
		}
	}(context.Background(), userID, posts)

	return posts, false, nil
}

type GetFeedResp []GetPostResp

func (s *Service) getFeed(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)

	posts, _, err := s.getPostsCached(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	resp := make(GetFeedResp, 0, len(posts))
	for _, post := range posts {
		resp = append(resp, GetPostResp{
			ID:     post.ID,
			Text:   post.Text,
			UserID: post.UserID,
		})
	}

	return c.JSON(http.StatusOK, resp)
}
