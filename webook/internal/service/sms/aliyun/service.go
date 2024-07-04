package aliyun

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"strings"
)

// AliSmsService 实现Service接口，用于发送阿里云短信
type AliSmsService struct {
	client   *dysmsapi.Client
	signName string
}

// NewAliSmsService 初始化AliSmsService实例
func NewAliSmsService(c *dysmsapi.Client, signName string) *AliSmsService {
	return &AliSmsService{
		client:   c,
		signName: signName,
	}
}

// Send 实现发送短信逻辑
func (s *AliSmsService) Send(ctx context.Context, tplId string, args map[string]string, numbers ...string) error {
	req := dysmsapi.CreateSendSmsRequest()
	req.Scheme = "https"
	// 阿里云多个手机为字符串逗号间隔
	req.PhoneNumbers = strings.Join(numbers, ",")
	req.SignName = s.signName
	bCode, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req.TemplateParam = string(bCode)
	req.TemplateCode = tplId
	var resp *dysmsapi.SendSmsResponse
	resp, err = s.client.SendSms(req)
	if err != nil {
		return err
	}

	if resp.Code != "OK" {
		return fmt.Errorf("短信发送失败，错误码：%s，错误信息：%s", resp.Code, resp.Message)
	}
	return nil
}
