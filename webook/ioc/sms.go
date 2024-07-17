package ioc

import (
	"GkWeiBook/webook/internal/service/sms"
	"GkWeiBook/webook/internal/service/sms/localsms"
)

func InitSMSService() sms.Service {
	return localsms.NewService()
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter())
	// 如果有需要，就可以用这个
	//return initTencentSMSService()
}

// 根据不同实现更换
//func initTencentSMSService() sms.Service {
//	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
//	if !ok {
//		panic("找不到腾讯 SMS 的 secret id")
//	}
//	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
//	if !ok {
//		panic("找不到腾讯 SMS 的 secret key")
//	}
//	c, err := tencentSMS.NewClient(
//		common.NewCredential(secretId, secretKey),
//		"ap-nanjing",
//		profile.NewClientProfile(),
//	)
//	if err != nil {
//		panic(err)
//	}
//	return tencent.NewService(c, "1400842696", "妙影科技")
//}
