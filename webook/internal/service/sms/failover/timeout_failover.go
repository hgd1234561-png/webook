package failover

import (
	"GkWeiBook/webook/internal/service/sms"
	"context"

	"sync/atomic"
)

// 连续N个超时就切换服务商

type TimeoutFailoverSMSService struct {
	svcs []sms.Service

	idx uint64

	// 连续超时的个数
	cnt uint64

	// 阈值; 连续超过这个数字，就要切换
	threshold uint64
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold uint64) sms.Service {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args map[string]string, numbers ...string) error {
	idx := atomic.LoadUint64(&t.idx)
	cnt := atomic.LoadUint64(&t.cnt)

	if cnt > t.threshold {
		// 这里要切换，新的下标，往后挪了一下
		newIdx := (idx + 1) % uint64(len(t.svcs))
		if atomic.CompareAndSwapUint64(&t.idx, idx, newIdx) {
			// 切换成功
			atomic.StoreUint64(&t.cnt, 0)
		}

		// else 就是出现并发，别人切换成功了
		idx = atomic.LoadUint64(&t.idx)
	}

	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddUint64(&t.cnt, 1)
	case nil:
		// 连续状态被打断
		atomic.StoreUint64(&t.cnt, 0)
	default:
		// 不知道什么错误
		return err
	}

	return err
}
