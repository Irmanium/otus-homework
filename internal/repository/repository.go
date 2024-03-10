package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"otus-homework/internal/domain"
)

type Repo struct {
	*pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{Conn: conn}
}

func (r *Repo) Register(ctx context.Context, user domain.FullUser) error {
	_, err := r.Exec(ctx, registerQuery, user.ID, user.FirstName, user.SecondName, user.Birthdate, user.Biography, user.City, user.PasswordHash)
	return err
}

func (r *Repo) GetUser(ctx context.Context, id string) (domain.UserProfile, error) {
	row := r.QueryRow(ctx, getUserQuery, id)

	var user domain.UserProfile
	err := row.Scan(&user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.UserProfile{}, domain.ErrNotFound
	}

	return user, err
}

func (r *Repo) GetPassword(ctx context.Context, id string) (string, error) {
	row := r.QueryRow(ctx, getPasswordQuery, id)

	var passwordHash string
	err := row.Scan(&passwordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrNotFound
	}

	return passwordHash, err
}
