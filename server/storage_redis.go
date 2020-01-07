package server

import (
	"context"
	"github.com/go-redis/redis/v7"
	"github.com/urfave/cli"
)

type redisStorage struct {
	client *redis.Client
}

var _ storage = &redisStorage{}

func newRedisStorage(c *cli.Context) (storage, error) {
	opt, err := redis.ParseURL(c.String("storage-args"))
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)
	return &redisStorage{client: client}, nil
}

func (l redisStorage) RecordLocation(ctx context.Context, state, location string) error {
	_, err := l.client.WithContext(ctx).Set(state, location, 0).Result()
	return err
}

func (l redisStorage) GetLocation(ctx context.Context, state string) (string, error) {
	pipe := l.client.WithContext(ctx).TxPipeline()
	getLoc := pipe.Get(state)
	pipe.Del(state)

	_, err := pipe.Exec()
	if err != nil {
		if err == redis.Nil {
			return "", errLocationNotFound
		}
		return "", err
	}

	location, err := getLoc.Result()
	if err != nil {
		if err == redis.Nil {
			return "", errLocationNotFound
		}
		return "", err
	}
	return location, nil
}

func (l redisStorage) Close() error {
	return l.client.Close()
}
