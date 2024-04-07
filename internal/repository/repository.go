package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"otus-homework/internal/domain"
)

type Repo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{Pool: conn}
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

	user.ID = id
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

func (r *Repo) SearchUser(ctx context.Context, firstName, secondName string) ([]domain.UserProfile, error) {
	rows, err := r.Query(ctx, searchUserQuery, firstName, secondName)
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
