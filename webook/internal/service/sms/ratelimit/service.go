package ratelimit

import (
	"GkWeiBook/webook/internal/service/sms"
	"GkWeiBook/webook/pkg/ratelimit"
	"context"
	"fmt"
)

// 利用装饰器进行 第三方服务治理

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSService) Send(ctx context.Context, tpl string, args map[string]string, numbers ...string) error {
	// 在这里加一些代码，新特性
	limited, err := s.limiter.Limit(ctx, "sms:localsms")
	if err != nil {
		// 系统错误，限流器出现问题
		// 可以限流：保守策略，你下游很坑
		// 可以不限，下游很强，保证更高的可用性
		return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)
	}
	if limited {
		// 限流了
		return fmt.Errorf("短信服务被限流了")
	}
	return s.svc.Send(ctx, tpl, args, numbers...)
}
