package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"otus-homework/internal/migrate"
	"otus-homework/internal/repo"
	"otus-homework/internal/service"
	"otus-homework/internal/shardedrepo"
)

const (
	dbString      = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"
	dbSlaveString = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"
	dbCitusString = "host=citus-master user=postgres password=postgres dbname=postgres sslmode=disable"

	repoMigrationsDir        = "migrations/repo"
	shardedRepoMigrationsDir = "migrations/sharded-repo"

	port          = "8080"
	jwtSecret     = "secret"
	tokenTTLHours = 72
)

func main() {
	migrate.Up(dbString, repoMigrationsDir)
	migrate.Up(dbCitusString, shardedRepoMigrationsDir)

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
	mainRepo := repo.New(pool, slavePool)

	shardedPool, err := pgxpool.New(context.Background(), dbCitusString)
	if err != nil {
		panic(err)
	}
	shardedRepo := shardedrepo.New(shardedPool)

	s := service.New(mainRepo, shardedRepo, port, jwtSecret, tokenTTLHours)
	s.StartService()
}
