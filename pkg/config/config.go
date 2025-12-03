package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 保存应用的所有配置设置
type Config struct {
	// 服务器配置
	ServerAddress string
	ServerPort    int

	// 数据库配置
	DatabaseDSN string

	// 运行环境
	Environment string
}

// NewConfig 创建并初始化一个新的 Config 实例
// 它会从环境变量以及可选的 .env 文件中加载配置
func NewConfig() (*Config, error) {
	// 如果 .env 文件存在则尝试加载（文件不存在时忽略错误）
	_ = godotenv.Load()

	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "0.0.0.0"),
		ServerPort:    getEnvAsInt("SERVER_PORT", 8080),
		DatabaseDSN:   getEnv("DATABASE_DSN", ""),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}

	// 校验必需的配置项
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// validate 检查所有必需配置项是否已设置
func (c *Config) validate() error {
	if c.DatabaseDSN == "" {
		return fmt.Errorf("DATABASE_DSN is required but not set")
	}
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("SERVER_PORT must be between 1 and 65535, got: %d", c.ServerPort)
	}
	return nil
}

// GetServerAddr 返回 host:port 格式的完整服务器地址
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.ServerAddress, c.ServerPort)
}

// IsProduction 判断应用是否运行在生产环境
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnv 获取环境变量的值，否则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量的整数值，否则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
