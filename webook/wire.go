//go:build wireinject

package main

import (
	"GkWeiBook/webook/internal/repository"
	"GkWeiBook/webook/internal/repository/cache"
	"GkWeiBook/webook/internal/repository/dao"
	"GkWeiBook/webook/internal/service"
	"GkWeiBook/webook/internal/web"
	"GkWeiBook/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 第三方依赖
		ioc.InitRedis, ioc.InitDB,
		// DAO 部分
		dao.NewUserDao,

		// cache 部分
		cache.NewCodeCache, cache.NewUserCache,

		// repository 部分
		repository.NewUserRepository,
		repository.NewCodeRepository,

		// Service 部分
		ioc.InitSMSService,
		service.NewUserService,
		service.NewCodeService,

		// handler 部分
		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
