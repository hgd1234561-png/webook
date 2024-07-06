package service

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/repository"
	repomocks "GkWeiBook/webook/internal/repository/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_userService_Login(t *testing.T) {

	tests := []struct {
		name string
		mack func(ctrl *gomock.Controller) repository.UserRepository

		email    string
		password string

		wantErr  error
		wantBody domain.User
	}{
		// TODO: Add test cases.
		{
			name: "登录成功",
			mack: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email: "123@qq.com",
					// 你在这边拿到的密码，就应该是一个正确的密码
					// 加密后的正确的密码
					Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
					Phone:    "15212345678",
				}, nil)
				return repo
			},

			email:    "123@qq.com",
			password: "123456#hello",
			wantBody: domain.User{
				Email: "123@qq.com",
				// 你在这边拿到的密码，就应该是一个正确的密码
				// 加密后的正确的密码
				Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
				Phone:    "15212345678",
			},
		},
		{
			name: "用户未找到",
			mack: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",
			wantErr:  ErrInvalidUserOrPassword,
		},

		{
			name: "系统错误",
			mack: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db错误"))
				return repo
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#hello",
			wantErr:  errors.New("db错误"),
		},

		{
			name: "密码不对",
			mack: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().
					FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email: "123@qq.com",
						// 你在这边拿到的密码，就应该是一个正确的密码
						// 加密后的正确的密码
						Password: "$2a$10$.l0JHmM7a2PdJ.A9gsmVyerEDlp1WhxsglC34S4UJH4TuHhWY7Tfq",
						Phone:    "15212345678",
					}, nil)
				return repo
			},
			email: "123@qq.com",
			// 用户输入的，没有加密的
			password: "123456#helloABCde",

			wantErr: ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tc.mack(ctrl)
			svc := NewUserService(repo)
			login, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantBody, login)
		})
	}
}
