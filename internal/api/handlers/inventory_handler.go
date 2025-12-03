package handlers

import (
	"bytes"
	"io"
	"net/http"
	"wms/internal/api/dto"
	"wms/internal/service"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InventoryHandler 负责处理与库存相关的 HTTP 请求
type InventoryHandler struct {
	service service.InventoryService
	logger  *logger.Logger
}

// NewInventoryHandler 创建一个新的 InventoryHandler 实例
func NewInventoryHandler(service service.InventoryService, log *logger.Logger) *InventoryHandler {
	return &InventoryHandler{
		service: service,
		logger:  log,
	}
}

// UploadCheck 处理库存盘点上传接口
// @Summary 上传库存盘点记录
// @Description 处理单次库存盘点操作
// @Tags inventory
// @Accept json
// @Produce json
// @Param request body dto.InventoryCheckRequest true "盘点请求数据"
// @Success 200 {object} dto.CommonResponse
// @Failure 400 {object} dto.CommonResponse "请求参数无效"
// @Failure 500 {object} dto.CommonResponse "服务器内部错误"
// @Router /api/wms/inventory/check/upload [post]
func (h *InventoryHandler) UploadCheck(c *gin.Context) {
	var req dto.InventoryCheckRequest

	// 记录原始请求体（用于调试）
	rawData, _ := c.GetRawData()
	h.logger.Debug("Received raw request body",
		zap.ByteString("raw_body", rawData),
		zap.String("remote_addr", c.ClientIP()),
	)

	// 重新设置请求体供后续绑定使用
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	// 绑定并校验 JSON 请求
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request payload",
			zap.Error(err),
			zap.ByteString("raw_body", rawData),
			zap.String("remote_addr", c.ClientIP()),
		)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse("Invalid request: "+err.Error()))
		return
	}

	h.logger.Info("Received inventory check request",
		zap.String("checker_id", req.CheckerID),
		zap.String("location_code", req.LocationCode),
		zap.String("material_code", req.MaterialCode),
		zap.Int("actual_quantity", req.ActualQuantity),
		zap.String("remote_addr", c.ClientIP()),
	)

	// 将 DTO 转换为服务层输入
	serviceInput := service.InventoryCheckInput{
		CheckerID:      req.CheckerID,
		LocationCode:   req.LocationCode,
		MaterialCode:   req.MaterialCode,
		ActualQuantity: req.ActualQuantity,
	}

	// 处理库存盘点
	if err := h.service.ProcessInventoryCheck(serviceInput); err != nil {
		h.logger.Error("Failed to process inventory check",
			zap.String("checker_id", req.CheckerID),
			zap.String("location_code", req.LocationCode),
			zap.String("material_code", req.MaterialCode),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse("Failed to process inventory check: "+err.Error()))
		return
	}

	h.logger.Info("Inventory check processed successfully",
		zap.String("checker_id", req.CheckerID),
		zap.String("material_code", req.MaterialCode),
		zap.String("location_code", req.LocationCode),
	)

	c.JSON(http.StatusOK, dto.SuccessResponse())
}
