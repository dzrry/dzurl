package redis

import (
	"context"
	"fmt"
	"github.com/dzrry/dzurl/domain"
	"github.com/dzrry/dzurl/service"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type redisRepo struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRepo(addr, port, password string) (*redisRepo, error) {
	repo := &redisRepo{}
	client, err := newClient(addr, port, password)
	if err != nil {
		return nil, fmt.Errorf("repository.Redis.NewRepo: %w", err)
	}
	repo.client = client
	return repo, nil
}

func newClient(addr, port, password string) (*redis.Client, error) {
	addr = fmt.Sprintf("%s:%s", addr, port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("repository.Redis.newClient: %w", err)
	}
	return client, nil
}

func (r *redisRepo) Load(key string) (*domain.Redirect, error) {
	redirect := &domain.Redirect{}
	data, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("repository.Redirect.Load: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("repository.Redirect.Load: %w", service.ErrRedirectNotFound)
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("repository.Redirect.Load: %w", err)
	}
	redirect.Key = data["key"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, nil
}

func (r *redisRepo) Store(redirect *domain.Redirect) error {
	key := redirect.Key
	data := map[string]interface{}{
		"key":        redirect.Key,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HSet(ctx, key, data).Result()
	if err != nil {
		return fmt.Errorf("repository.Redirect.Store: %w", err)
	}
	return nil
}
