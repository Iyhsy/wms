package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// globalLogger 是默认的日志实例
	globalLogger *Logger
)

// Logger 包装 zap.Logger，提供结构化日志并区分业务与技术日志
type Logger struct {
	zapLogger *zap.Logger
	sugar     *zap.SugaredLogger
}

// NewLogger 根据指定环境创建新的 Logger 实例
// 生产环境：输出 INFO 级别及以上的 JSON 日志
// 开发环境：输出易读的控制台日志并启用 DEBUG 级别
func NewLogger(environment string) (*Logger, error) {
	var zapLogger *zap.Logger
	var err error

	if environment == "production" {
		// 生产环境：JSON 输出并使用 INFO 等级
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zapLogger, err = config.Build(
			zap.AddCallerSkip(1), // 在堆栈追踪中跳过包装函数
			zap.AddStacktrace(zap.ErrorLevel),
		)
	} else {
		// 开发环境：控制台输出并使用 DEBUG 等级
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapLogger, err = config.Build(
			zap.AddCallerSkip(1),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	logger := &Logger{
		zapLogger: zapLogger,
		sugar:     zapLogger.Sugar(),
	}

	return logger, nil
}

// InitGlobalLogger 初始化全局日志实例
func InitGlobalLogger(environment string) error {
	logger, err := NewLogger(environment)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetLogger 返回全局日志实例
// 如果尚未初始化，则创建默认的开发环境日志实例
func GetLogger() *Logger {
	if globalLogger == nil {
		logger, err := NewLogger("development")
		if err != nil {
			// 如果创建日志失败则回退至 stderr
			fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
			return nil
		}
		globalLogger = logger
	}
	return globalLogger
}

// Info 以 INFO 等级输出包含结构化字段的日志
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

// Warn 以 WARN 等级输出包含结构化字段的日志
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

// Error 以 ERROR 等级输出包含结构化字段的日志
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, fields...)
}

// Debug 以 DEBUG 等级输出包含结构化字段的日志
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

// Fatal 以 FATAL 等级输出日志并退出应用
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, fields...)
}

// Infof 以 INFO 等级输出格式化消息（便捷方法）
func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warnf 以 WARN 等级输出格式化消息（便捷方法）
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Errorf 以 ERROR 等级输出格式化消息（便捷方法）
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Debugf 以 DEBUG 等级输出格式化消息（便捷方法）
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// With 创建带额外字段的子 logger
func (l *Logger) With(fields ...zap.Field) *Logger {
	newZapLogger := l.zapLogger.With(fields...)
	return &Logger{
		zapLogger: newZapLogger,
		sugar:     newZapLogger.Sugar(),
	}
}

// Sync 刷新缓冲的日志条目
func (l *Logger) Sync() error {
	return l.zapLogger.Sync()
}

// 包级别快捷函数，依赖全局 logger

// Info 使用全局 logger 以 INFO 等级输出日志
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn 使用全局 logger 以 WARN 等级输出日志
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error 使用全局 logger 以 ERROR 等级输出日志
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Debug 使用全局 logger 以 DEBUG 等级输出日志
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Fatal 使用全局 logger 以 FATAL 等级输出日志并退出
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}
