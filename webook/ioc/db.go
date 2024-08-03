package ioc

import (
	"GkWeiBook/webook/internal/repository/dao"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	type config struct {
		DSN string `yaml:"dsn"`
	}
	var c config
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(fmt.Sprintf("配置文件解析失败: %v", err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
