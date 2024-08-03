package ioc

import (
	"GkWeiBook/webook/internal/service/oauth2/wechat"
	"GkWeiBook/webook/internal/web"
)

func InitWechatService() wechat.Service {
	//appId, ok := os.LookupEnv("WECHAT_APP_ID")
	//if !ok {
	//	panic("WECHAT_APP_ID not found")
	//}
	//appSecret, ok := os.LookupEnv("WECHAT_APP_SECRET")
	//if !ok {
	//	panic("WECHAT_APP_SECRET not found")
	//}
	//return wechat.NewService(appId, appSecret)
	return wechat.NewService("appId", "")

}

func InitWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
