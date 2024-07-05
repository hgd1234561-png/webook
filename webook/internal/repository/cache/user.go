package cache

import (
	"GkWeiBook/webook/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, du domain.User) error
}
type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 30,
	}
}

// error为nil 默认为有数据
// 如果没有数据 返回一个特点的error
func (cache *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := cache.key(uid)
	// 我假定这个地方用 JSON 来
	data, err := cache.client.Get(ctx, key).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(data), &u)
	//if err != nil {
	//	return domain.User{}, err
	//}
	//return u, nil
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, du domain.User) error {
	key := cache.key(du.Id)
	// 我假定这个地方用 JSON
	data, err := json.Marshal(du)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, key, data, cache.expiration).Err()
}

func (cache *RedisUserCache) key(uid int64) string {
	// user-info-
	// user.info.
	// user/info/
	// user_info_
	return fmt.Sprintf("user:info:%d", uid)
}
