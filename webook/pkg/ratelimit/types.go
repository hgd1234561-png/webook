package ratelimit

import "context"

// 整体限流 不要把第三方服务打崩

type Limiter interface {
	// Limited 是否触发限流，key 是限流对象
	// 如果触发限流，返回 true error 代表限流器是否有错误

	Limit(ctx context.Context, key string) (bool, error)
}
