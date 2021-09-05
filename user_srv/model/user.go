package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	Id        uint           `gorm:"primaryKey;AUTO_INCREMENT;unsigned;comment:id"`
	CreatedAt time.Time      `gorm:"column:add_time;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:update_time;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:";comment:删除时间"`
	IsDelete  bool           `gorm:";comment:是否删除"`
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"type:varchar(11);unique;index:idx_mobile;not null;comment:手机号"`
	Password string     `gorm:"type:varchar(100);not null;comment:密码"`
	NickName string     `gorm:"type:varchar(20);comment:昵称"`
	Birthday *time.Time `gorm:"type:datetime;comment:生日"`
	Gender   uint8      `gorm:"default:0;comment:0男1女"`
	Role     uint8      `gorm:"default:1; comment '1普通用户2管理员用户'"`
}
