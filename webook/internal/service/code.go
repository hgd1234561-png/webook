package service

import (
	"GkWeiBook/webook/internal/repository"
	"GkWeiBook/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

// 验证码服务
type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *CodeService) Send(ctx context.Context,
	// 业务标识，区别业务场景
	biz string,
	phone string) error {
	// 生成验证码
	code := svc.generateCode()
	//塞进redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	const codeTplId = "SMS_468730234" //换成自己的
	// 发送出去
	err = svc.smsSvc.Send(ctx, codeTplId, map[string]string{
		"code": code,
	}, phone)
	//if err != nil {
	//	// 意味着Redis存成功了，但是短信发送失败了
	//	//可以选择重试
	//	//但这里选择让用户重新发
	//
	//}

	return err
}

func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	ok, err := svc.repo.Verify(ctx, biz, phone, inputCode)
	if err == repository.ErrCodeVerifyTooMany {
		// 相当于，我们对外面屏蔽了验证次数过多的错误，我们就是告诉调用者，你这个不对
		return false, nil
	}
	return ok, err
}

func (svc *CodeService) generateCode() string {
	// 六位数，0-999999包括0和999999
	code := rand.Intn(1000000)
	//格式化一下
	return fmt.Sprintf("%06d", code)
}
