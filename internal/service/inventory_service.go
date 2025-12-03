package service

import (
	"fmt"
	"time"
	"wms/internal/model"
	"wms/internal/repository"
	"wms/pkg/logger"

	"go.uber.org/zap"
)

// InventoryCheckInput 表示盘点处理所需的输入数据
type InventoryCheckInput struct {
	CheckerID      string `json:"checker_id"`
	LocationCode   string `json:"location_code"`
	MaterialCode   string `json:"material_code"`
	ActualQuantity int    `json:"actual_quantity"`
}

// InventoryService 定义库存业务逻辑接口
type InventoryService interface {
	// ProcessInventoryCheck 处理单条库存盘点操作
	// 在一个事务中依次执行以下步骤：
	// 1. 获取指定库位的当前库存
	// 2. 计算差异（difference = actual_quantity - stock_quantity）
	// 3. 生成盘点记录
	// 4. 将库存数量更新为实盘数量
	ProcessInventoryCheck(input InventoryCheckInput) error

	// ProcessBatchInventoryCheck 处理多条盘点任务
	// 每条盘点都使用独立事务单独处理
	ProcessBatchInventoryCheck(inputs []InventoryCheckInput) []error
}

// inventoryService 是 InventoryService 的具体实现
type inventoryService struct {
	repo   repository.InventoryCheckRepository
	logger *logger.Logger
}

// NewInventoryService 创建一个新的 InventoryService 实例
func NewInventoryService(repo repository.InventoryCheckRepository, log *logger.Logger) InventoryService {
	return &inventoryService{
		repo:   repo,
		logger: log,
	}
}

// ProcessInventoryCheck 在事务安全前提下处理单条盘点操作
func (s *inventoryService) ProcessInventoryCheck(input InventoryCheckInput) error {
	// 输入校验
	if err := s.validateInput(input); err != nil {
		s.logger.Error("Invalid inventory check input",
			zap.String("checker_id", input.CheckerID),
			zap.String("location_code", input.LocationCode),
			zap.String("material_code", input.MaterialCode),
			zap.Error(err),
		)
		return fmt.Errorf("input validation failed: %w", err)
	}

	s.logger.Info("Starting inventory check process",
		zap.String("checker_id", input.CheckerID),
		zap.String("location_code", input.LocationCode),
		zap.String("material_code", input.MaterialCode),
		zap.Int("actual_quantity", input.ActualQuantity),
	)

	// 开启事务
	tx := s.repo.BeginTransaction()
	if tx.Error != nil {
		s.logger.Error("Failed to begin transaction",
			zap.String("material_code", input.MaterialCode),
			zap.String("location_code", input.LocationCode),
			zap.Error(tx.Error),
		)
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// 确保事务被正确处理
	defer func() {
		if r := recover(); r != nil {
			s.repo.RollbackTransaction(tx)
			s.logger.Error("Panic during inventory check, transaction rolled back",
				zap.String("material_code", input.MaterialCode),
				zap.Any("panic", r),
			)
		}
	}()

	// 查询当前库存
	stock, err := s.repo.GetStockByMaterialAndLocation(tx, input.MaterialCode, input.LocationCode)
	if err != nil {
		s.repo.RollbackTransaction(tx)
		s.logger.Error("Failed to fetch stock information",
			zap.String("material_code", input.MaterialCode),
			zap.String("location_code", input.LocationCode),
			zap.Error(err),
		)
		return fmt.Errorf("failed to fetch stock: %w", err)
	}

	// 若库存不存在则创建初始库存
	var stockQuantity int
	if stock == nil {
		s.logger.Info("Stock record not found, creating new stock entry",
			zap.String("material_code", input.MaterialCode),
			zap.String("location_code", input.LocationCode),
		)
		stockQuantity = 0
		stock = &model.Stock{
			MaterialCode: input.MaterialCode,
			LocationCode: input.LocationCode,
			Quantity:     input.ActualQuantity,
		}
	} else {
		stockQuantity = stock.Quantity
		stock.Quantity = input.ActualQuantity
	}

	// 计算差异
	difference := input.ActualQuantity - stockQuantity

	s.logger.Debug("Calculated inventory variance",
		zap.String("material_code", input.MaterialCode),
		zap.String("location_code", input.LocationCode),
		zap.Int("stock_quantity", stockQuantity),
		zap.Int("actual_quantity", input.ActualQuantity),
		zap.Int("difference", difference),
	)

	// 创建盘点记录
	checkRecord := &model.InventoryCheckRecord{
		CheckerID:      input.CheckerID,
		LocationCode:   input.LocationCode,
		MaterialCode:   input.MaterialCode,
		ActualQuantity: input.ActualQuantity,
		StockQuantity:  stockQuantity,
		Difference:     difference,
		CheckTime:      time.Now(),
		IsProcessed:    true,
	}

	if err := s.repo.CreateCheckRecord(tx, checkRecord); err != nil {
		s.repo.RollbackTransaction(tx)
		s.logger.Error("Failed to create inventory check record",
			zap.String("material_code", input.MaterialCode),
			zap.String("location_code", input.LocationCode),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create check record: %w", err)
	}

	// 更新库存数量
	if err := s.repo.UpdateStock(tx, stock); err != nil {
		s.repo.RollbackTransaction(tx)
		s.logger.Error("Failed to update stock quantity",
			zap.String("material_code", input.MaterialCode),
			zap.String("location_code", input.LocationCode),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// 提交事务
	if err := s.repo.CommitTransaction(tx); err != nil {
		s.repo.RollbackTransaction(tx)
		s.logger.Error("Failed to commit transaction",
			zap.String("material_code", input.MaterialCode),
			zap.String("location_code", input.LocationCode),
			zap.Error(err),
		)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Inventory check completed successfully",
		zap.String("checker_id", input.CheckerID),
		zap.String("material_code", input.MaterialCode),
		zap.String("location_code", input.LocationCode),
		zap.Int("previous_quantity", stockQuantity),
		zap.Int("new_quantity", input.ActualQuantity),
		zap.Int("variance", difference),
	)

	return nil
}

// ProcessBatchInventoryCheck 独立处理多条盘点任务
func (s *inventoryService) ProcessBatchInventoryCheck(inputs []InventoryCheckInput) []error {
	s.logger.Info("Starting batch inventory check process",
		zap.Int("total_items", len(inputs)),
	)

	errors := make([]error, len(inputs))
	successCount := 0
	failureCount := 0

	for i, input := range inputs {
		err := s.ProcessInventoryCheck(input)
		errors[i] = err
		if err != nil {
			failureCount++
		} else {
			successCount++
		}
	}

	s.logger.Info("Batch inventory check completed",
		zap.Int("total_items", len(inputs)),
		zap.Int("success_count", successCount),
		zap.Int("failure_count", failureCount),
	)

	return errors
}

// validateInput 校验盘点输入参数
func (s *inventoryService) validateInput(input InventoryCheckInput) error {
	if input.CheckerID == "" {
		return fmt.Errorf("checker_id is required")
	}
	if input.LocationCode == "" {
		return fmt.Errorf("location_code is required")
	}
	if input.MaterialCode == "" {
		return fmt.Errorf("material_code is required")
	}
	if input.ActualQuantity < 0 {
		return fmt.Errorf("actual_quantity cannot be negative: %d", input.ActualQuantity)
	}
	return nil
}
