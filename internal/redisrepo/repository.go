package redisrepo

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"otus-homework/internal/domain"
)

type Repo struct {
	*redis.Client
}

func New(client *redis.Client) *Repo {
	return &Repo{
		Client: client,
	}
}

func getUserFeedKey(userID string) string {
	return "feed_" + userID
}

func (r *Repo) PutFeedToCache(ctx context.Context, userID string, feed []domain.Post) error {
	val, err := json.Marshal(&feed)
	if err != nil {
		return err
	}

	return r.Set(ctx, getUserFeedKey(userID), val, 0).Err()
}

func (r *Repo) GetFeedFromCache(ctx context.Context, userID string) ([]domain.Post, error) {
	val, err := r.Get(ctx, getUserFeedKey(userID)).Result()
	if err != nil {
		return nil, err
	}

	var feed []domain.Post
	err = json.Unmarshal([]byte(val), &feed)

	return feed, err
}
