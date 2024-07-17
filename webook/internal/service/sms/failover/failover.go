package failover

import (
	"GkWeiBook/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

// 用装饰器模式实现短信服务商的切换

type FailOverSMSService struct {
	svcs []sms.Service

	// 当前服务商下标
	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) sms.Service {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

//func (f FailOverSMSService) Send(ctx context.Context, tpl string, args map[string]string, numbers ...string) error {
//	// 最简单的轮询
//	for _, svc := range f.svcs {
//		err := svc.Send(ctx, tpl, args, numbers...)
//		if err == nil {
//			return nil
//		}
//		log.Println("发送短信失败，切换下一个服务商")
//	}
//
//	return errors.New("所有服务商都发送失败")
//}

func (f *FailOverSMSService) Send(ctx context.Context, tpl string, args map[string]string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	// 我要迭代 length
	for i := idx; i < idx+length; i++ {
		// 取余数来计算下标
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// 前者是被取消，后者是超时
			return err
		}
		log.Println(err)
	}
	return errors.New("轮询了所有的服务商，但是发送都失败了")
}
