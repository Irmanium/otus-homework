package service

import (
	"context"
	"time"

	"otus-homework/internal/domain"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserRepo interface {
	Register(ctx context.Context, user domain.FullUser) error
	GetUser(ctx context.Context, id string) (domain.UserProfile, error)
	GetPassword(ctx context.Context, id string) (string, error)
	SearchUser(ctx context.Context, firstName, secondName string) ([]domain.UserProfile, error)

	GetFriends(ctx context.Context, userID string) ([]string, error)
	SetFriend(ctx context.Context, ids [2]string) error
	DeleteFriend(ctx context.Context, ids [2]string) error

	CreatePost(ctx context.Context, id, userID, text string) error
	UpdatePost(ctx context.Context, id, text string) error
	DeletePost(ctx context.Context, id string) error
	GetPost(ctx context.Context, id string) (userID, text string, updatedAt time.Time, err error)

	GetFeed(ctx context.Context, userID string) ([]domain.Post, error)
}

type DialogRepo interface {
	SendMessage(ctx context.Context, id string, message domain.Message) error
	GetDialog(ctx context.Context, userID, interlocutorID string) ([]domain.Message, error)
}

type FeedCache interface {
	PutFeedToCache(ctx context.Context, userID string, feed []domain.Post) error
	GetFeedFromCache(ctx context.Context, userID string) ([]domain.Post, error)
}

type LiveFeedRepo interface {
	SendPost(ctx context.Context, post domain.Post) error
	GetFeed(friendIDs []string, cancel <-chan struct{}) (<-chan []byte, error)
}

type Service struct {
	*echo.Echo

	userRepo     UserRepo
	dialogRepo   DialogRepo
	feedCache    FeedCache
	liveFeedRepo LiveFeedRepo

	port          string
	jwtSecret     string
	tokenTTLHours int

	feedCacheMaxLen int
}

type ErrResp struct {
	Message string `json:"message"`
}

func New(
	userRepo UserRepo,
	dialogRepo DialogRepo,
	feedCache FeedCache,
	liveFeedRepo LiveFeedRepo,
	port string,
	jwtSecret string,
	tokenTTLHours int,
	feedCacheMaxLen int,
) *Service {

	e := echo.New()

	return &Service{
		Echo:            e,
		userRepo:        userRepo,
		dialogRepo:      dialogRepo,
		feedCache:       feedCache,
		liveFeedRepo:    liveFeedRepo,
		port:            port,
		jwtSecret:       jwtSecret,
		tokenTTLHours:   tokenTTLHours,
		feedCacheMaxLen: feedCacheMaxLen,
	}
}

func (s *Service) StartService() {
	s.Use(middleware.Logger())
	s.Use(middleware.Recover())

	s.POST("/login", s.login)
	s.POST("/user/register", s.register)

	r := s.Group("")
	r.Use(echojwt.WithConfig(s.getJWTConfig()))

	r.GET("/user/get/:id", s.getUser)
	r.GET("/user/search", s.searchUser)

	r.POST("/dialog/:user_id/send", s.sendMessage)
	r.GET("/dialog/:user_id/list", s.getDialog)

	r.PUT("/friend/set/:user_id", s.setFriend)
	r.PUT("/friend/delete/:user_id", s.deleteFriend)

	r.POST("/post/create", s.createPost)
	r.PUT("/post/update", s.updatePost)
	r.PUT("/post/delete/:id", s.deletePost)
	r.GET("/post/get/:id", s.getPost)
	r.GET("/post/feed", s.getFeed)
	r.GET("/post/feed/posted", s.postFeed) // ws

	s.Logger.Fatal(s.Start(":" + s.port))
}

func convertDomainProfileToResp(user domain.UserProfile) GetUserResp {
	return GetUserResp{
		ID:         user.ID,
		FirstName:  user.FirstName,
		SecondName: user.SecondName,
		Birthdate:  user.Birthdate.Format(time.DateOnly),
		Biography:  user.Biography,
		City:       user.City,
	}
}

func (s *Service) rebuildUserFeedCache(userID string) {
	go func(ctx context.Context, userID string) {
		posts, err := s.userRepo.GetFeed(ctx, userID)
		if err != nil {
			s.Logger.Warn(err)
			return
		}

		err = s.feedCache.PutFeedToCache(ctx, userID, posts)
		if err != nil {
			s.Logger.Warn(err)
		}
	}(context.Background(), userID)
}
