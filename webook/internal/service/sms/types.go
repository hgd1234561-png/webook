package sms

import "context"

// Service 发送短信的抽象
// 屏蔽不同供应商之间的区别
type Service interface {
	Send(ctx context.Context, tpl string, args map[string]string, numbers ...string) error
}

type TencentService interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}
