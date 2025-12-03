package model

import "time"

// InventoryCheckRecord 表示单条库存盘点记录
// 该结构保存实盘数量与系统库存的对比信息
type InventoryCheckRecord struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CheckerID      string    `gorm:"type:varchar(100);not null;index" json:"checker_id"`
	LocationCode   string    `gorm:"type:varchar(100);not null;index" json:"location_code"`
	MaterialCode   string    `gorm:"type:varchar(100);not null;index" json:"material_code"`
	ActualQuantity int       `gorm:"not null" json:"actual_quantity"`
	StockQuantity  int       `gorm:"not null" json:"stock_quantity"`
	Difference     int       `gorm:"not null" json:"difference"`
	CheckTime      time.Time `gorm:"type:timestamp;not null;index" json:"check_time"`
	IsProcessed    bool      `gorm:"default:false;not null" json:"is_processed"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定 InventoryCheckRecord 对应的表名
func (InventoryCheckRecord) TableName() string {
	return "inventory_check_records"
}
