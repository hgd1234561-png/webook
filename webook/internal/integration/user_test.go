package integration

import (
	"GkWeiBook/webook/ioc"
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 集成测试

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	test := []struct {
		name string

		before func(t *testing.T)
		after  func(t *testing.T)

		phone string

		wantCode int
		wantBody string
	}{
		{
			name: "发送成功的用例",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				code, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				assert.NoError(t, err)
				// 你的验证码是六位
				assert.True(t, len(code) == 6)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: `"` + "发送成功" + `"`,
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := rdb.Set(ctx, "phone_code:login:15212345678", "123456", time.Minute*9+time.Second*50).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				code, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				assert.NoError(t, err)
				// 你的验证码是六位
				assert.Equal(t, "123456", code)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: `"` + "短信发送太频繁，请稍后再试" + `"`,
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				err := rdb.Set(ctx, "phone_code:login:15212345678", "123456", 0).Err()
				assert.NoError(t, err)

			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				defer cancel()
				code, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				assert.NoError(t, err)
				// 你的验证码是六位
				assert.Equal(t, "123456", code)
			},
			phone:    "15212345678",
			wantCode: http.StatusOK,
			wantBody: `"` + "系统错误" + `"`,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			tt.before(t)
			// 准备Req和记录的 recorder
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone": "%s"}`, tt.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tt.wantCode, recorder.Code)
			assert.Equal(t, tt.wantBody, recorder.Body.String())

			tt.after(t)
		})
	}
}
