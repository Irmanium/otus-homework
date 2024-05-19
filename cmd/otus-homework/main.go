package main

import (
	"context"
	"otus-homework/internal/tarantoolrepo"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tarantool/go-tarantool/v2"
	"otus-homework/internal/migrate"
	"otus-homework/internal/repo"
	"otus-homework/internal/service"
)

const (
	tarantoolAddress = "tarantool:3301"

	dbString      = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"
	dbSlaveString = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"

	repoMigrationsDir = "migrations/repo"

	port          = "8080"
	jwtSecret     = "secret"
	tokenTTLHours = 72
)

func main() {
	tarantoolConn, err := tarantool.Connect(context.Background(), tarantool.NetDialer{
		Address: tarantoolAddress,
	}, tarantool.Opts{})
	if err != nil {
		panic(err)
	}
	tarantoolRepo, err := tarantoolrepo.NewAndMigrate(tarantoolConn)
	if err != nil {
		panic(err)
	}

	migrate.Up(dbString, repoMigrationsDir)
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

	s := service.New(mainRepo, tarantoolRepo, port, jwtSecret, tokenTTLHours)
	s.StartService()
}
