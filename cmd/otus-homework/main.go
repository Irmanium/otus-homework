package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"otus-homework/internal/migrate"
	"otus-homework/internal/repository"
	"otus-homework/internal/service"
)

const (
	dbString      = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"
	port          = "8080"
	jwtSecret     = "secret"
	tokenTTLHours = 72
)

func main() {
	migrate.Up(dbString)

	conn, err := pgx.Connect(context.Background(), dbString)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = conn.Close(context.Background())
		if err != nil {
			panic(err)
		}
	}()
	repo := repository.New(conn)

	s := service.New(repo, port, jwtSecret, tokenTTLHours)
	s.StartService()
}
