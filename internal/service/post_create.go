package service

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"otus-homework/internal/domain"
)

func (s *Service) addPostInUserFeedCache(userID, postID, text, authorID string) {
	go func(ctx context.Context, userID, postID, text string) {
		posts, fromCache, err := s.getPostsCached(ctx, userID)
		if err != nil {
			s.Logger.Warn(err)
			return
		}

		if !fromCache {
			return
		}

		posts = append(posts, domain.Post{
			ID:     postID,
			Text:   text,
			UserID: authorID,
		})
		if len(posts) > s.feedCacheMaxLen {
			posts = posts[0:s.feedCacheMaxLen]
		}

		err = s.feedCache.PutFeedToCache(ctx, userID, posts)
		if err != nil {
			s.Logger.Warn(err)
		}
	}(context.Background(), userID, postID, text)
}

func (s *Service) addPostInFriendsFeedCache(userID, postID, text string) {
	go func(ctx context.Context, userID string) {
		friends, err := s.userRepo.GetFriends(ctx, userID)
		if err != nil {
			s.Logger.Warn(err)
			return
		}

		for _, friend := range friends {
			s.addPostInUserFeedCache(friend, postID, text, userID)
		}
	}(context.Background(), userID)
}

type CreatePostReq struct {
	Text string `json:"text"`
}

type CreatePostResp struct {
	ID string `json:"id"`
}

func (s *Service) createPost(c echo.Context) error {
	userID := s.getUserIDFromJWT(c)

	req := new(CreatePostReq)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrResp{Message: "invalid request"})
	}

	id := uuid.New().String()

	if err := s.userRepo.CreatePost(c.Request().Context(), id, userID, req.Text); err != nil {
		return err
	}

	s.addPostInFriendsFeedCache(userID, id, req.Text)

	return c.JSON(http.StatusOK, CreatePostResp{ID: id})
}
