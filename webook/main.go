package main

import (
	"GkWeiBook/webook/internal/repository"
	"GkWeiBook/webook/internal/repository/dao"
	"GkWeiBook/webook/internal/service"
	"GkWeiBook/webook/internal/web"
	"GkWeiBook/webook/internal/web/middleware"
	"GkWeiBook/webook/pkg/ginx/middlewares/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

// 前端node版本是18.16.0
func main() {

	//db := initDB()
	//
	//u := initUser(db)
	//
	//r := initWebServer()
	//
	//u.RegisterRouters(r)
	//
	//r.Run(":8080")
	r := gin.Default()

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(200, "Hello, K8S!")
	})
	r.GET("/du", func(ctx *gin.Context) {
		ctx.String(200, "傻ber,你还真看啊！")
	})

	r.Run(":8080")
}

func initWebServer() *gin.Engine {
	r := gin.Default()

	// 跨域
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
		//AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 是否允许带cookie之类的东西
		AllowCredentials: true,
		// 不加这个前端是拿不到的
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.Contains(origin, "http://localhost") {
				return true
			}
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	// session
	//store := cookie.NewStore([]byte("secret"))
	//store, err := redis.NewStore(16, "tcp", "8.130.82.184:6379", "",
	//	[]byte("bstmmTdM2KFxXcm544kMZzzBsBgwgb6J"), []byte("xYGqV38Qbz4VVvbYJdQqK5XnEVkP6SPe"))
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//r.Use(sessions.Sessions("mysession", store))

	redisClient := redis.NewClient(&redis.Options{
		Addr: "8.130.82.184:6379",
	})
	r.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	r.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").IgnorePaths("/users/login").Build())

	return r
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(8.130.82.184:3306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
