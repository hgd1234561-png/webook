package service

import (
	"GkWeiBook/webook/internal/domain"
	"GkWeiBook/webook/internal/repository"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var ErrCodeSendTooMany = repository.ErrCodeSendTooMany
var ErrUserDuplicateEmail = repository.ErrUserDuplicateUser
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或者密码错误")

type UserService interface {
	Profile(ctx context.Context, userId int64) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放在哪里的问题
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 然后就是存起来
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 先去查一下用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 然后就是校验密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *userService) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.UpdateUserById(ctx, u)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先找一下，我们认为，大部分用户是已经存在的用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// 有两种情况
		// err == nil, u 是可用的
		// err != nil，系统错误，
		return u, err
	}
	// 用户没找到
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	// 有两种可能，一种是 err 恰好是唯一索引冲突（phone）
	// 一种是 err != nil，系统错误
	if err != nil && err != repository.ErrUserDuplicateUser {
		return domain.User{}, err
	}
	// 要么 err ==nil，要么ErrDuplicateUser，也代表用户存在
	// 主从延迟，理论上来讲，强制走主库
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenId)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 这边就是意味着是一个新用户
	err = svc.repo.Create(ctx, domain.User{
		WechatInfo: info,
	})
	if err != nil && err != repository.ErrUserDuplicateUser {
		return domain.User{}, err
	}
	return svc.repo.FindByWechat(ctx, info.OpenId)
}

func (svc *userService) Profile(ctx context.Context, userId int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userId)
}
