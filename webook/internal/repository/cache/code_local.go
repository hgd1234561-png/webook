package cache

import (
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/net/context"
	"sync"
	"time"
)

// 用本地缓存代替Redis

// 技术选型考虑的点
// 1. 功能性: 功能是否能够完全覆盖你的需求。
// 2. 社区和支持度: 社区是否活跃, 文档是否齐全, 以及搜索引擎能不能搜索到各种信息
// 3. 非功能性: 易用性（用户友好度，学习曲线要平滑）,扩展性,性能

type LocalCodeCache struct {
	cache      *lru.Cache
	lock       sync.Mutex
	expiration time.Duration
}

func NewLocalCodeCache(c *lru.Cache, expiration time.Duration) CodeCache {

	return &LocalCodeCache{
		cache:      c,
		expiration: expiration,
	}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		// 说明没有验证码
		l.cache.Add(key, codeItem{
			code:      code,
			cnt:       3,
			expiredAt: now.Add(l.expiration),
		})
		return nil
	}
	itm, ok := val.(codeItem)
	if !ok {
		// 理论上不太可能
		return errors.New("类型转换失败")
	}
	if itm.expiredAt.Sub(now) > time.Minute*9 {
		// 不到一分钟
		return ErrCodeSendTooMany
	}

	// 重新发送
	l.cache.Add(key, codeItem{
		code:      code,
		cnt:       3,
		expiredAt: now.Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	val, ok := l.cache.Get(key)
	if !ok {
		// 说明没有验证码
		return false, ErrKeyNotExist
	}
	itm, ok := val.(codeItem)
	if !ok {
		// 理论上不太可能
		return false, errors.New("类型转换失败")
	}
	if itm.cnt <= 0 {
		return false, ErrCodeVerifyTooMany
	}
	itm.cnt--
	return itm.code == inputCode, nil
}

func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItem struct {
	code      string
	cnt       int
	expiredAt time.Time
}
