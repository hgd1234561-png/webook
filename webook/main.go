package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"
)

// 前端node版本是18.16.0
func main() {

	//initViper()

	initViperRemote()
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

func initViper() {
	cfile := pflag.String("dev",
		"webook/config/dev.yaml", "配置文件路径")
	// 这一步之后，cfile 里面才有值
	pflag.Parse()
	//viper.Set("db.dsn", "localhost:3306")
	// 所有的默认值放好s
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3",
		"http://8.130.82.184:2379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	log.Println("远程配置中心发生变更")
	//})
	//go func() {
	//	for {
	//		err = viper.WatchRemoteConfig()
	//		if err != nil {
	//			panic(err)
	//		}
	//		log.Println("watch", viper.GetString("test.key"))
	//		//time.Sleep(time.Second)
	//	}
	//}()
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
