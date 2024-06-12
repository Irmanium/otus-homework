package repo

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"otus-homework/internal/domain"
)

type Repo struct {
	*pgxpool.Pool
	slave *pgxpool.Pool

	feedMaxLen int
}

func New(conn *pgxpool.Pool, slaveConn *pgxpool.Pool, feedMaxLen int) *Repo {
	return &Repo{
		Pool:       conn,
		slave:      slaveConn,
		feedMaxLen: feedMaxLen,
	}
}

func (r *Repo) Register(ctx context.Context, user domain.FullUser) error {
	_, err := r.Exec(ctx, registerQuery, user.ID, user.FirstName, user.SecondName, user.Birthdate, user.Biography, user.City, user.PasswordHash)
	return err
}

func (r *Repo) GetUser(ctx context.Context, id string) (domain.UserProfile, error) {
	row := r.slave.QueryRow(ctx, getUserQuery, id)

	var user domain.UserProfile
	err := row.Scan(&user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.UserProfile{}, domain.ErrNotFound
	}

	user.ID = id
	return user, err
}

func (r *Repo) GetPassword(ctx context.Context, id string) (string, error) {
	row := r.slave.QueryRow(ctx, getPasswordQuery, id)

	var passwordHash string
	err := row.Scan(&passwordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrNotFound
	}

	return passwordHash, err
}

func (r *Repo) SearchUser(ctx context.Context, firstName, secondName string) ([]domain.UserProfile, error) {
	rows, err := r.slave.Query(ctx, searchUserQuery, firstName, secondName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.UserProfile
	for rows.Next() {
		var user domain.UserProfile
		err = rows.Scan(&user.ID, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City)
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repo) SetFriend(ctx context.Context, ids [2]string) error {
	_, err := r.Exec(ctx, setFriendQuery, ids[0], ids[1])
	return err
}

func (r *Repo) DeleteFriend(ctx context.Context, ids [2]string) error {
	_, err := r.Exec(ctx, deleteFriendQuery, ids[0], ids[1])
	return err
}

func (r *Repo) CreatePost(ctx context.Context, id, userID, text string) error {
	_, err := r.Exec(ctx, createPostQuery, id, userID, text)
	return err
}

func (r *Repo) UpdatePost(ctx context.Context, id, text string) error {
	_, err := r.Exec(ctx, updatePostQuery, id, text)
	return err
}

func (r *Repo) DeletePost(ctx context.Context, id string) error {
	_, err := r.Exec(ctx, deletePostQuery, id)
	return err
}

func (r *Repo) GetPost(ctx context.Context, id string) (userID, text string, updatedAt time.Time, err error) {
	row := r.slave.QueryRow(ctx, getPostQuery, id)

	err = row.Scan(&userID, &text, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", time.Time{}, domain.ErrNotFound
	}

	return userID, text, updatedAt, err
}

func (r *Repo) GetFeed(ctx context.Context, userID string) ([]domain.Post, error) {
	rows, err := r.slave.Query(ctx, getFeedQuery, userID, r.feedMaxLen)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		err = rows.Scan(&post.ID, &post.Text, &post.UserID)
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repo) GetFriends(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.slave.Query(ctx, getFriendsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
