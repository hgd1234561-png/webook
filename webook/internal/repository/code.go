package repository

import (
	"GkWeiBook/webook/internal/repository/cache"
	"context"
)

var ErrCodeSendTooMany = cache.ErrCodeSendTooMany
var ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCachedCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

func (repo *CachedCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CachedCodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
