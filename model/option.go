package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Option 系统配置实体
type Option struct {
	ID              uint      `gorm:"primarykey" json:"id"`
	Key             string    `gorm:"column:key;size:100;uniqueIndex;not null" json:"key"`
	Value           string    `gorm:"type:text;not null" json:"value"`
	GroupType       string    `gorm:"size:50;not null;default:system;index" json:"groupType"`
	Component       string    `gorm:"size:50" json:"component"`                     // 前端组件：input, textarea, select, switch, date-picker 等
	ComponentParams JSONValue `gorm:"type:json" json:"componentParams"`             // 组件参数
	CreatedAt       time.Time `gorm:"not null;index" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"not null" json:"updatedAt"`
}

// TableName 指定表名
func (Option) TableName() string {
	return "options"
}

// JSONValue JSON值类型
type JSONValue map[string]interface{}

// Scan 实现 sql.Scanner 接口
func (j *JSONValue) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口
func (j JSONValue) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// String 转为JSON字符串
func (j JSONValue) String() string {
	if j == nil {
		return ""
	}
	bytes, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(bytes)
}
