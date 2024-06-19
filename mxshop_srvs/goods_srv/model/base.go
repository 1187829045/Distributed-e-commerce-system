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

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}

type BaseModel struct {
	ID        int32          `gorm:"primarykey;type:int" json:"id"` // 主键，类型为 int32，对应数据库中的 int 类型
	CreatedAt time.Time      `gorm:"column:add_time" json:"-"`      // 创建时间，对应数据库中的 add_time 列，JSON 序列化时忽略此字段
	UpdatedAt time.Time      `gorm:"column:update_time" json:"-"`   // 更新时间，对应数据库中的 update_time 列，JSON 序列化时忽略此字段
	DeletedAt gorm.DeletedAt `json:"-"`                             // 软删除时间，gorm.DeletedAt 是 gorm 框架提供的软删除支持类型，JSON 序列化时忽略此字段
	IsDeleted bool           `json:"-"`                             // 是否已删除，JSON 序列化时忽略此字段
}
