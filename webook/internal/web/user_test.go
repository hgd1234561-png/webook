package web

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/service"
	svcmocks "GkWeiBook/webook/internal/service/mocks"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string

		// mock
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)

		// 构造请求，预期中输入
		reqBuilder func(t *testing.T) *http.Request

		// 预期中的输出
		wantCode int
		wantBody string
	}{
		// TODO: Add test cases.
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "邮箱格式不正确",
		},
		{
			name: "bind错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"passsword": "hello#world123",
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusBadRequest,
		},
		{
			name: "两次密码输入不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "helloworld123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "两次密码输入不一致",
		},
		{
			name: "密码不符合要求",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello",
"confirmPassword": "hello"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位, 包含数字、特殊字符",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(errors.New("db错误"))
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(service.ErrUserDuplicateEmail)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/signup", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123",
"confirmPassword": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造 handler
			userSvc, codeSvc := tt.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			// 准备服务器，注册路由
			server := gin.Default()
			hdl.RegisterRouters(server)

			// 准备Req和记录的 recorder
			req := tt.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			// 断言结果
			assert.Equal(t, tt.wantCode, recorder.Code)
			assert.Equal(t, tt.wantBody, recorder.Body.String())
		})
	}
}

func TestUserHandler_LoginJWT(t *testing.T) {

	tests := []struct {
		name string

		// mock
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		// 构造请求，预期中输入
		reqBuilder func(t *testing.T) *http.Request

		wantbody string
		wantcode int
	}{
		// TODO: Add test cases.
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "123@qq.com", "hello#world123").Return(domain.User{
					Id: int64(1),
				}, nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/login", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantcode: http.StatusOK,
			wantbody: "登录成功",
		},
		{
			name: "用户名或密码错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "123@qq.com", "hello#world123").Return(domain.User{}, service.ErrInvalidUserOrPassword)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/login", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantcode: http.StatusOK,
			wantbody: "用户名或密码错误",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "123@qq.com", "hello#world123").Return(domain.User{}, errors.New("db错误"))
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost,
					"/users/login", bytes.NewReader([]byte(`{
"email": "123@qq.com",
"password": "hello#world123"
}`)))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},

			wantcode: http.StatusOK,
			wantbody: "系统错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 构造 handler
			userSvc, codeSvc := tt.mock(ctrl)

			hdl := NewUserHandler(userSvc, codeSvc)

			// 准备服务器，注册路由
			server := gin.Default()
			hdl.RegisterRouters(server)

			// 准备Req和记录的 recorder
			req := tt.reqBuilder(t)
			recorder := httptest.NewRecorder()

			// 执行
			server.ServeHTTP(recorder, req)

			fmt.Println(" ")
			assert.Equal(t, tt.wantcode, recorder.Code)
			assert.Equal(t, tt.wantbody, recorder.Body.String())
		})
	}
}
