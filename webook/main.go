package main

import (
	"github.com/gin-gonic/gin"
)

// 前端node版本是18.16.0
func main() {

	r := InitWebServer()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, K8S!")
	})
	r.GET("/du", func(ctx *gin.Context) {
		ctx.String(200, "傻ber,你还真看啊！")
	})
	r.Run(":8080")

}

//func initWebServer() *gin.Engine {
//	r := gin.Default()
//
//	// 跨域
//	r.Use(cors.New(cors.Config{
//		AllowOrigins: []string{"http://localhost:3000"},
//		//AllowMethods: []string{"POST", "GET"},
//		AllowHeaders: []string{"Content-Type", "Authorization"},
//		// 是否允许带cookie之类的东西
//		AllowCredentials: true,
//		// 不加这个前端是拿不到的
//		ExposeHeaders: []string{"x-jwt-token"},
//		AllowOriginFunc: func(origin string) bool {
//			if strings.Contains(origin, "http://localhost") {
//				return true
//			}
//			return true
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//
//	// session
//	//store := cookie.NewStore([]byte("secret"))
//	//store, err := redis.NewStore(16, "tcp", "8.130.82.184:6379", "",
//	//	[]byte("bstmmTdM2KFxXcm544kMZzzBsBgwgb6J"), []byte("xYGqV38Qbz4VVvbYJdQqK5XnEVkP6SPe"))
//	//
//	//if err != nil {
//	//	panic(err)
//	//}
//	//
//	//r.Use(sessions.Sessions("mysession", store))
//
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: "101.126.22.227:30399",
//	})
//	r.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
//
//	r.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/login_sms/code/send").
//		IgnorePaths("/users/login_sms").
//		IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())
//
//	return r
//}
//
//func initUser(db *gorm.DB) *web.UserHandler {
//	ud := dao.NewUserDao(db)
//	rs := redis.NewClient(&redis.Options{
//		//Addr: "webook-record-redis:6379",
//		Addr: "101.126.22.227:30399",
//	})
//	uc := cache.NewUserCache(rs)
//	repo := repository.NewUserRepository(ud, uc)
//	svc := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(rs)
//	codePo := repository.NewCodeRepository(codeCache)
//	codeSvc := service.NewCodeService(codePo, localsms.NewService())
//	u := web.NewUserHandler(svc, codeSvc)
//	return u
//}
//
//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open("root:root@tcp(101.126.22.227:31823)/webook"))
//	if err != nil {
//		panic(err)
//	}
//
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//	return db
//}
