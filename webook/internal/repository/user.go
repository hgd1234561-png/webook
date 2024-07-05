package repository

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/repository/cache"
	"GkWeiBook/webook/internal/repository/dao"
	"context"
	"database/sql"
	"time"
)

var (
	ErrUserDuplicateUser = dao.ErrUserDuplicateEmail
	ErrUserNotFound      = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUserById(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.toEntity(u))
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *CachedUserRepository) UpdateUserById(ctx context.Context, u domain.User) error {

	return r.dao.Update(ctx, r.toEntity(u))
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		// 必然有数据
		return u, nil
	}
	//// 没有数据
	//if err == cache.ErrKeyNotExist {
	//	// 去数据库加载
	//}
	// err 不为 nil，就要查询数据库
	// err 有两种可能
	// 1. key 不存在，说明 redis 是正常的
	// 2. 访问 redis 有问题。可能是网络有问题，也可能是 redis 本身就崩溃了
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.toDomain(ue)

	err = r.cache.Set(ctx, u)
	if err != nil {
		//打日志做监控
	}
	return u, err
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.toDomain(u), nil
}

func (r *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		Password: u.Password,
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
		Birthday: time.UnixMilli(u.Birthday),
	}
}

func (r *CachedUserRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Birthday: u.Birthday.UnixMilli(),
		AboutMe:  u.AboutMe,
		Nickname: u.Nickname,
	}
}
