package repository

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/repository/cache"
	cachemocks "GkWeiBook/webook/internal/repository/cache/mocks"
	"GkWeiBook/webook/internal/repository/dao"
	daomocks "GkWeiBook/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {

	nowMs := time.Now().UnixMilli()
	//now := time.UnixMilli(nowMs)

	tests := []struct {
		name string

		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)
		id   int64

		want    domain.User
		wantErr error
	}{
		// TODO: Add test cases.
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id: uid,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "123456",
					Birthday: 100,
					AboutMe:  "自我介绍",
					Phone: sql.NullString{
						String: "15212345678",
						Valid:  true,
					},
					Ctime: nowMs,
					Utime: 102,
				}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "15212345678",
					Password: "123456",
					Birthday: time.UnixMilli(100),
					AboutMe:  "自我介绍",
				}).Return(nil)

				return d, c
			},

			id: 123,
			want: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "15212345678",
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "15212345678",
					Password: "123456",
					Birthday: time.UnixMilli(100),
					AboutMe:  "自我介绍",
				}, nil)

				return d, c
			},

			id: 123,
			want: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "15212345678",
			},
			wantErr: nil,
		},
		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{}, dao.ErrUserNotFound)

				return d, c
			},

			id:      123,
			want:    domain.User{},
			wantErr: dao.ErrUserNotFound,
		},
		{
			name: "回写缓存失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				uid := int64(123)
				d := daomocks.NewMockUserDAO(ctrl)
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id: uid,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "123456",
					Birthday: 100,
					AboutMe:  "自我介绍",
					Phone: sql.NullString{
						String: "15212345678",
						Valid:  true,
					},
					Ctime: nowMs,
					Utime: 102,
				}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Phone:    "15212345678",
					Password: "123456",
					Birthday: time.UnixMilli(100),
					AboutMe:  "自我介绍",
				}).Return(errors.New("redis错误"))

				return d, c
			},

			id: 123,
			want: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Birthday: time.UnixMilli(100),
				AboutMe:  "自我介绍",
				Phone:    "15212345678",
			},
			wantErr: errors.New("redis错误"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, ca := tt.mock(ctrl)
			repo := NewCachedUserRepository(ud, ca)
			user, err := repo.FindById(context.Background(), tt.id)
			assert.Equal(t, tt.want, user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
