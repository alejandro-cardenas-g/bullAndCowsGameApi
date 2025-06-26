package store

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr, pw string, db int) (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})

	if err := r.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return r, nil
}
