package db

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB() (*pgxpool.Pool, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL not set")
	}
	dbpool,err:=pgxpool.New(context.Background(), databaseURL)
	if err!=nil{
		return nil,err
	}

	return dbpool,nil
}
