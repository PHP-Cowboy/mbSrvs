package main

import (
	"crypto/sha512"
	"fmt"
	"log"
	"mbSrvs/user_srv/model"
	"os"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

//func genMd5(code string) string {
//	md := md5.New()
//	_, _ = io.WriteString(md, code)
//	return hex.EncodeToString(md.Sum(nil))
//}

func main() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	logger := logger2.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger2.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			Colorful:      false,       //禁用彩色打印
			LogLevel:      logger2.Info,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix: "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: logger,
	})
	if err != nil {
		panic(err)
	}

	//_ = db.AutoMigrate(&model.User{})

	birthday, err := time.ParseInLocation("2006-01-02 15:04:04", "1998-07-26 00:00:00", time.Local)

	if err != nil {
		panic("时间转换出错:" + err.Error())
	}

	for i := 0; i < 10; i++ {
		user := model.User{
			Mobile:   fmt.Sprintf("1570018000%d", i),
			Password: GenPwd(fmt.Sprintf("admin@#%d", i)),
			NickName: fmt.Sprintf("小王子%d", i),
			Birthday: &birthday,
			Gender:   0,
			Role:     1,
		}
		db.Save(&user)
	}

	//db.Create(&user)

}

func GenPwd(pwd string) string {
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, genPwd := password.Encode(pwd, options)
	//check := password.Verify("123456",salt,pwd,options)

	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, genPwd)
}
