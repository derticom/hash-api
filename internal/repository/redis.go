package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/derticom/hash-api/internal/domain"
)

type RedisRepository struct {
	client *redis.Client
	logger *slog.Logger
	ttl    time.Duration
}

// NewRepository создает новый экземпляр Redis репозитория.
func NewRepository(client *redis.Client, log *slog.Logger, ttl time.Duration) *RedisRepository {
	return &RedisRepository{
		client: client,
		logger: log,
		ttl:    ttl,
	}
}

// Store сохраняет данные хеша в кеш по исходной строке.
func (r *RedisRepository) Store(ctx context.Context, input string, data *domain.HashData) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		r.logger.Error("Failed to marshal data",
			slog.String("input", input),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to marshal hash data: %w", err)
	}

	r.client.Set(ctx, input, jsonData, r.ttl)

	r.logger.Debug("Input data stored successfully",
		slog.String("input", input),
		slog.String("md5", data.MD5),
		slog.String("sha256", data.SHA256),
	)

	return nil
}

// GetByInput получает данные из кеша по исходной строке.
func (r *RedisRepository) GetByInput(ctx context.Context, input string) (*domain.HashData, error) {
	jsonData, err := r.client.Get(ctx, input).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.logger.Debug("Input not found in cache",
				slog.String("input", input),
			)
			return nil, domain.ErrHashNotFound
		}

		return nil, fmt.Errorf("failed to client.Get: %w", err)
	}

	var data domain.HashData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to json.Unmarshal: %w", err)
	}

	r.logger.Debug("Input data retrieved successfully",
		slog.String("input", input),
		slog.String("md5", data.MD5),
		slog.String("sha256", data.SHA256),
	)

	return &data, nil
}
