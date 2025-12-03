package model

import "time"

// Stock 表示特定库位中某物料的当前库存
type Stock struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	MaterialCode string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_material_location" json:"material_code"`
	LocationCode string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_material_location" json:"location_code"`
	Quantity     int       `gorm:"not null;default:0" json:"quantity"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定 Stock 对应的表名
func (Stock) TableName() string {
	return "stocks"
}
