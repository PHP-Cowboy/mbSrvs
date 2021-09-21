package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

func (g *GormList) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), g)
}

type BaseModel struct {
	Id        int32          `gorm:"primaryKey;AUTO_INCREMENT;unsigned;comment:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:add_time;comment:创建时间" json:"-"`
	UpdatedAt time.Time      `gorm:"column:update_time;comment:更新时间" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"comment:删除时间" json:"-"`
	IsDelete  bool           `gorm:"comment:是否删除" json:"-"`
}
