package main

import (
	"context"

	"otus-homework/internal/migrate"
	"otus-homework/internal/rabbitrepo"
	"otus-homework/internal/redisrepo"
	"otus-homework/internal/repo"
	"otus-homework/internal/service"
	"otus-homework/internal/tarantoolrepo"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/tarantool/go-tarantool/v2"
)

const (
	tarantoolAddress = "tarantool:3301"

	dbString      = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"
	dbSlaveString = "host=db user=postgres password=postgres dbname=postgres sslmode=disable"

	redisAddress = "redis:6379"

	rabbitMQString = "amqp://guest:guest@rabbitmq:5672/"

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

	conn, err := amqp.Dial(rabbitMQString)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	rabbitRepo, err := rabbitrepo.New(conn)
	if err != nil {
		panic(err)
	}

	s := service.New(mainRepo, tarantoolRepo, redisRepo, rabbitRepo, port, jwtSecret, tokenTTLHours, feedMaxLen)
	s.StartService()
}
