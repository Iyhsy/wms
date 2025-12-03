package repository

import (
	"wms/internal/model"

	"gorm.io/gorm"
)

// InventoryCheckRepository 定义库存盘点数据访问接口
type InventoryCheckRepository interface {
	// CreateCheckRecord 在事务中创建新的盘点记录
	CreateCheckRecord(tx *gorm.DB, record *model.InventoryCheckRecord) error

	// GetStockByMaterialAndLocation 查询指定库位与物料的库存信息
	GetStockByMaterialAndLocation(tx *gorm.DB, materialCode, locationCode string) (*model.Stock, error)

	// UpdateStock 更新指定库位物料的库存数量
	UpdateStock(tx *gorm.DB, stock *model.Stock) error

	// BeginTransaction 开启新的数据库事务
	BeginTransaction() *gorm.DB

	// CommitTransaction 提交当前事务
	CommitTransaction(tx *gorm.DB) error

	// RollbackTransaction 回滚当前事务
	RollbackTransaction(tx *gorm.DB) error

	// FindCheckRecordsByMaterial 查询某物料的全部盘点记录
	FindCheckRecordsByMaterial(materialCode string) ([]model.InventoryCheckRecord, error)

	// FindUnprocessedRecords 查询所有未处理的盘点记录
	FindUnprocessedRecords() ([]model.InventoryCheckRecord, error)
}

// inventoryCheckRepository 是 InventoryCheckRepository 的具体实现
type inventoryCheckRepository struct {
	db *gorm.DB
}

// NewInventoryCheckRepository 创建新的 InventoryCheckRepository 实例
func NewInventoryCheckRepository(db *gorm.DB) InventoryCheckRepository {
	return &inventoryCheckRepository{
		db: db,
	}
}

// CreateCheckRecord 在事务中新增盘点记录
func (r *inventoryCheckRepository) CreateCheckRecord(tx *gorm.DB, record *model.InventoryCheckRecord) error {
	return tx.Create(record).Error
}

// GetStockByMaterialAndLocation 查询指定库位的库存
// 若库存不存在则返回 nil（不视为错误）
func (r *inventoryCheckRepository) GetStockByMaterialAndLocation(tx *gorm.DB, materialCode, locationCode string) (*model.Stock, error) {
	var stock model.Stock
	err := tx.Where("material_code = ? AND location_code = ?", materialCode, locationCode).First(&stock).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

// UpdateStock 更新指定库位的库存数量
func (r *inventoryCheckRepository) UpdateStock(tx *gorm.DB, stock *model.Stock) error {
	return tx.Save(stock).Error
}

// BeginTransaction 开启新的数据库事务
func (r *inventoryCheckRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

// CommitTransaction 提交当前事务
func (r *inventoryCheckRepository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

// RollbackTransaction 回滚当前事务
func (r *inventoryCheckRepository) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

// FindCheckRecordsByMaterial 查询指定物料的所有盘点记录
func (r *inventoryCheckRepository) FindCheckRecordsByMaterial(materialCode string) ([]model.InventoryCheckRecord, error) {
	var records []model.InventoryCheckRecord
	err := r.db.Where("material_code = ?", materialCode).Order("check_time DESC").Find(&records).Error
	return records, err
}

// FindUnprocessedRecords 查询所有未处理的盘点记录
func (r *inventoryCheckRepository) FindUnprocessedRecords() ([]model.InventoryCheckRecord, error) {
	var records []model.InventoryCheckRecord
	err := r.db.Where("is_processed = ?", false).Order("check_time ASC").Find(&records).Error
	return records, err
}
