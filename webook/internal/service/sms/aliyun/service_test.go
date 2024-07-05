package aliyun

import (
	"context"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"os"
	"testing"
)

func TestSender(t *testing.T) {
	config := sdk.NewConfig()

	credential := credentials.NewStsTokenCredential(os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_ID"), os.Getenv("ALIBABA_CLOUD_ACCESS_KEY_SECRET"), os.Getenv("ALIBABA_CLOUD_SECURITY_TOKEN"))

	client, err := dysmsapi.NewClientWithOptions("cn-wulanchabu", config, credential)
	if err != nil {
		panic(err)
	}

	aliSmsService := NewAliSmsService(client, "webook")
	err = aliSmsService.Send(context.Background(), "SMS_468730234", map[string]string{
		"code": "1234",
	}, "")
	if err != nil {
		panic(err)
	}

}
