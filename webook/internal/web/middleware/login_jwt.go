package middleware

import (
	ijwt "GkWeiBook/webook/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenStr := l.ExtractToken(ctx)
		uc := ijwt.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return []byte("bstmmTdM2KFxXcm544kMZzzBsBgwgb6J"), nil
		})

		if err != nil {
			// 没登陆,有人瞎搞 Bearer xxxxxx
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// err 为 nil, token 不为 nil
		if !token.Valid {
			// 没登陆
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != ctx.Request.UserAgent() {
			//严重的安全问题
			//你是要监控的
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		err = l.CheckSession(ctx, uc.Ssid)
		if err != nil {
			// redis有问题， 或者退出登录了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user", uc)
	}
}
