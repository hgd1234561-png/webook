package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ErrCodeSendTooMany = errors.New("验证码发送太频繁")
var ErrCodeVerifyTooMany = errors.New("验证码校验太频繁")

//编译器会在编译的时候，把set_code的代码放进来这个luaSetCode变量里

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		//没有问题
		return nil
	case -1:
		//发送太频繁
		return ErrCodeSendTooMany
	default:
		return errors.New("系统错误")

	}

}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		// 调用 redis 出了问题
		return false, err
	}
	switch res {
	case -2:
		return false, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	default:
		return true, nil
	}
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
