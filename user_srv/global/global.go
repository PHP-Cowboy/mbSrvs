package global

import (
	"fmt"
	"log"
	"mbSrvs/user_srv/config"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB           *gorm.DB
	ServerConfig = &config.ServerConfig{}
)

func init() {

	info := ServerConfig.MysqlInfo

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		info.User,
		info.Password,
		info.Host,
		info.Port,
		info.Name,
	)

	logger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      false,       //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}
}
