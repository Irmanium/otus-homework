package main

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/tarantool/go-tarantool/v2"
	"otus-homework/internal/migrate"
	"otus-homework/internal/redisrepo"
	"otus-homework/internal/repo"
	"otus-homework/internal/service"
	"otus-homework/internal/tarantoolrepo"
)

const (
	tarantoolAddress = "tarantool:3301"

	dbString      = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"
	dbSlaveString = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"

	redisAddress = "redis:6379"

	repoMigrationsDir = "migrations/repo"

	port          = "8080"
	jwtSecret     = "secret"
	tokenTTLHours = 72
	feedMaxLen    = 1000
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
	mainRepo := repo.New(pool, slavePool, feedMaxLen)

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	redisRepo := redisrepo.New(redisClient)

	s := service.New(mainRepo, tarantoolRepo, redisRepo, port, jwtSecret, tokenTTLHours, feedMaxLen)
	s.StartService()
}
