package ioc

import (
	"GkWeiBook/webook/internal/web"
	ijwt "GkWeiBook/webook/internal/web/jwt"
	"GkWeiBook/webook/internal/web/middleware"
	"GkWeiBook/webook/pkg/logger"
	"golang.org/x/net/context"

	rlimit "GkWeiBook/webook/pkg/ginx/middlewares/ratelimit"
	"GkWeiBook/webook/pkg/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, oauth2WechatHandler *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRouters(server)
	oauth2WechatHandler.RegisterRouters(server)
	return server
}

func InitGinMiddlewares(limiter ratelimit.Limiter, jwtHdl ijwt.Handler, log logger.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowAllOrigins: true,
			//AllowOrigins:     []string{"http://localhost:3000"},
			AllowCredentials: true,

			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 这个是允许前端访问你的后端响应中带的头部
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
			//AllowHeaders: []string{"content-type"},
			//AllowMethods: []string{"POST"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					//if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "your_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		rlimit.NewBuilder(limiter).Build(),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/users/signup").IgnorePaths("/users/login").Build(),
		middleware.NewLogMiddlewareBuilder(func(ctx context.Context, l middleware.AccessLog) {
			log.Debug("", logger.Field{Key: "req", Val: l})
		}).AllowReqBody().AllowRespBody().Build(),
	}
}
