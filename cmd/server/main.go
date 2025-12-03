package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wms/internal/api/handlers"
	"wms/internal/api/routes"
	"wms/internal/model"
	"wms/internal/repository"
	"wms/internal/service"
	"wms/pkg/config"
	"wms/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log, err := logger.NewLogger(cfg.Environment)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	log.Info("Starting WMS Inventory System",
		zap.String("environment", cfg.Environment),
		zap.String("server_address", cfg.GetServerAddr()),
	)

	// 建立数据库连接
	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database",
			zap.Error(err),
			zap.String("dsn", cfg.DatabaseDSN),
		)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("Database connection established successfully")

	// 自动迁移数据库模型
	if err := db.AutoMigrate(&model.Stock{}, &model.InventoryCheckRecord{}); err != nil {
		log.Fatal("Failed to auto-migrate database models", zap.Error(err))
	}

	log.Info("Database migration completed successfully")

	// 手动注入依赖
	// 仓储层
	inventoryRepo := repository.NewInventoryCheckRepository(db)

	// 服务层
	inventoryService := service.NewInventoryService(inventoryRepo, log)

	// 处理器层
	inventoryHandler := handlers.NewInventoryHandler(inventoryService, log)

	// 初始化 Gin 路由
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// 注册 API 路由
	routes.SetupRoutes(router, inventoryHandler)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: router,
	}

	// 在 goroutine 中启动服务器
	go func() {
		log.Info("HTTP server starting",
			zap.String("address", cfg.GetServerAddr()),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	log.Info("WMS Inventory System is running",
		zap.String("address", cfg.GetServerAddr()),
		zap.String("environment", cfg.Environment),
	)

	// 优雅停机机制
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// 阻塞直至接收到信号
	sig := <-quit
	log.Info("Shutdown signal received",
		zap.String("signal", sig.String()),
	)

	// 创建 5 秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试执行优雅停机
	log.Info("Initiating graceful shutdown", zap.Duration("timeout", 5*time.Second))

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	// 关闭数据库连接
	if err := sqlDB.Close(); err != nil {
		log.Error("Failed to close database connections", zap.Error(err))
	} else {
		log.Info("Database connections closed successfully")
	}

	log.Info("WMS Inventory System shutdown complete")
}
