package global

import (
	"gorm.io/gorm"
	"mbSrvs/user_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig = &config.ServerConfig{}
)
