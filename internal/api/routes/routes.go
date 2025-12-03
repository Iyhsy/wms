package routes

import (
	"wms/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 配置应用的所有路由
func SetupRoutes(router *gin.Engine, inventoryHandler *handlers.InventoryHandler) {
	// API v1 路由分组
	api := router.Group("/api/wms")
	{
		// 库存相关路由
		inventory := api.Group("/inventory")
		{
			check := inventory.Group("/check")
			{
				check.POST("/upload", inventoryHandler.UploadCheck)
			}
		}
	}
}
