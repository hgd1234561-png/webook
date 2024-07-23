package ratelimit

import (
	"GkWeiBook/webook/internal/service/sms"
	smsmocks "GkWeiBook/webook/internal/service/sms/mocks"
	"GkWeiBook/webook/pkg/ratelimit"
	"GkWeiBook/webook/pkg/ratelimitmocks"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"testing"
)

func TestRatelimitSMSService_Send(t *testing.T) {

	tests := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter)

		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "限流成功",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := ratelimitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(context.Background(), "sms:localsms").Return(true, nil)
				return svc, limiter
			},
			wantErr: errors.New("短信服务被限流了"),
		},

		{
			name: "正常发送",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := ratelimitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(context.Background(), "sms:localsms").Return(false, nil)
				svc.EXPECT().Send(context.Background(), "mysql", map[string]string{
					"code": "1234",
				}, "12345678901").Return(nil)
				return svc, limiter
			},
			wantErr: nil,
		},

		{
			name: "限流器异常",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := ratelimitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(context.Background(), "sms:localsms").Return(false, errors.New("限流器异常"))
				return svc, limiter
			},
			wantErr: fmt.Errorf("短信服务判断是否限流出现问题，%w", errors.New("限流器异常")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc, limiter := tt.mock(ctrl)
			limitSvc := NewRatelimitSMSService(svc, limiter)
			err := limitSvc.Send(context.Background(), "mysql", map[string]string{
				"code": "1234",
			}, "12345678901")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
