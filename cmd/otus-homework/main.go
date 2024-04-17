package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"otus-homework/internal/migrate"
	"otus-homework/internal/repository"
	"otus-homework/internal/service"
)

const (
	dbString      = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"
	dbSlaveString = "host=db-slave-one user=postgres password=postgres dbname=postgres sslmode=disable"
	port          = "8080"
	jwtSecret     = "secret"
	tokenTTLHours = 72
)

func main() {
	migrate.Up(dbString)

	pool, err := pgxpool.New(context.Background(), dbString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	slavePool, err := pgxpool.New(context.Background(), dbSlaveString)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	repo := repository.New(pool, slavePool)

	s := service.New(repo, port, jwtSecret, tokenTTLHours)
	s.StartService()
}

