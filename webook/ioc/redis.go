package ioc

import (
	"GkWeiBook/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

func InitRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
	})
}

func InitLimiter(cmd redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSlidingWindowLimiter(cmd, time.Second, 1000)
}
