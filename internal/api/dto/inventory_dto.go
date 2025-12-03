package dto

// InventoryCheckRequest 表示库存盘点上传的请求负载
type InventoryCheckRequest struct {
	CheckerID      string `json:"checker_id" binding:"required"`
	LocationCode   string `json:"location_code" binding:"required"`
	MaterialCode   string `json:"material_code" binding:"required"`
	ActualQuantity int    `json:"actual_quantity" binding:"required,min=0"`
}

// CommonResponse 表示标准的 API 响应结构
type CommonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SuccessResponse 创建成功响应
func SuccessResponse() CommonResponse {
	return CommonResponse{
		Code:    0,
		Message: "success",
	}
}

// ErrorResponse 创建带自定义消息的错误响应
func ErrorResponse(message string) CommonResponse {
	return CommonResponse{
		Code:    -1,
		Message: message,
	}
}
