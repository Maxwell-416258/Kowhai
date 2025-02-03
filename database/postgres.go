package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"kowhai/global"
	"log"
)

func InitPostgres() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		global.Config.Postgres.User, global.Config.Postgres.Password, global.Config.Postgres.Host, global.Config.Postgres.Port, global.Config.Postgres.Database))
	if err != nil {
		log.Fatalf("无法连接 PostgreSQL: %v", err)
		return nil, err
	}
	return dbpool, nil
}
