package repository_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/derticom/hash-api/internal/domain"
	"github.com/derticom/hash-api/internal/repository"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisRepository_StoreAndGet(t *testing.T) {
	s, err := miniredis.Run()
	require.NoError(t, err)
	defer s.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	repo := repository.NewRepository(rdb, logger, 10*time.Minute)

	ctx := context.Background()
	input := "test_input"
	hashData := &domain.HashData{
		MD5:    "md5_value",
		SHA256: "sha256_value",
	}

	// Store
	err = repo.Store(ctx, input, hashData)
	require.NoError(t, err)

	// Get
	result, err := repo.GetByInput(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, hashData.MD5, result.MD5)
	assert.Equal(t, hashData.SHA256, result.SHA256)

	// Пытаемся получить несуществующий ключ
	result, err = repo.GetByInput(ctx, "not_exist")
	assert.ErrorIs(t, err, domain.ErrHashNotFound)
	assert.Nil(t, result)
}
