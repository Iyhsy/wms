package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wms/internal/service"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
)

// mockInventoryService 是用于测试的模拟实现
type mockInventoryService struct {
	processFunc func(input service.InventoryCheckInput) error
}

func (m *mockInventoryService) ProcessInventoryCheck(input service.InventoryCheckInput) error {
	if m.processFunc != nil {
		return m.processFunc(input)
	}
	return nil
}

func (m *mockInventoryService) ProcessBatchInventoryCheck(inputs []service.InventoryCheckInput) []error {
	return nil
}

func setupTestRouter(handler *InventoryHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/wms/inventory/check/upload", handler.UploadCheck)
	return router
}

func TestUploadCheck_Success(t *testing.T) {
	// 准备环境
	log, _ := logger.NewLogger("test")
	mockService := &mockInventoryService{}
	handler := NewInventoryHandler(mockService, log)
	router := setupTestRouter(handler)

	// 测试数据
	requestBody := map[string]interface{}{
		"checker_id":      "user123",
		"location_code":   "LOC001",
		"material_code":   "MAT001",
		"actual_quantity": 100,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 构造请求
	req, _ := http.NewRequest("POST", "/api/wms/inventory/check/upload", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != 0 {
		t.Errorf("Expected code 0, got %v", response["code"])
	}
	if response["message"] != "success" {
		t.Errorf("Expected message 'success', got %v", response["message"])
	}
}

func TestUploadCheck_MissingRequiredField(t *testing.T) {
	// 准备环境
	log, _ := logger.NewLogger("test")
	mockService := &mockInventoryService{}
	handler := NewInventoryHandler(mockService, log)
	router := setupTestRouter(handler)

	// 测试数据——缺少 checker_id
	requestBody := map[string]interface{}{
		"location_code":   "LOC001",
		"material_code":   "MAT001",
		"actual_quantity": 100,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 构造请求
	req, _ := http.NewRequest("POST", "/api/wms/inventory/check/upload", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != -1 {
		t.Errorf("Expected code -1, got %v", response["code"])
	}
}

func TestUploadCheck_NegativeQuantity(t *testing.T) {
	// 准备环境
	log, _ := logger.NewLogger("test")
	mockService := &mockInventoryService{}
	handler := NewInventoryHandler(mockService, log)
	router := setupTestRouter(handler)

	// 测试数据——负数库存
	requestBody := map[string]interface{}{
		"checker_id":      "user123",
		"location_code":   "LOC001",
		"material_code":   "MAT001",
		"actual_quantity": -10,
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 构造请求
	req, _ := http.NewRequest("POST", "/api/wms/inventory/check/upload", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != -1 {
		t.Errorf("Expected code -1, got %v", response["code"])
	}
}

func TestUploadCheck_InvalidDataType(t *testing.T) {
	// 准备环境
	log, _ := logger.NewLogger("test")
	mockService := &mockInventoryService{}
	handler := NewInventoryHandler(mockService, log)
	router := setupTestRouter(handler)

	// 测试数据——actual_quantity 使用字符串而非整数
	requestBody := map[string]interface{}{
		"checker_id":      "user123",
		"location_code":   "LOC001",
		"material_code":   "MAT001",
		"actual_quantity": "invalid",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 构造请求
	req, _ := http.NewRequest("POST", "/api/wms/inventory/check/upload", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 断言
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["code"].(float64) != -1 {
		t.Errorf("Expected code -1, got %v", response["code"])
	}
}
