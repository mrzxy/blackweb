package lib

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gormLogger "gorm.io/gorm/logger"
)

var (
	Logger *zap.SugaredLogger
)

func InitLogger() {
	// 初始化zap日志
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = "time"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, err := cfg.Build()
	if err != nil {
		panic("初始化zap失败: " + err.Error())
	}
	zap.ReplaceGlobals(logger)
	Logger = logger.Sugar()
}

type DBLogger struct {
	logger *zap.SugaredLogger
}

func NewDBLogger(logger *zap.SugaredLogger) *DBLogger {
	return &DBLogger{
		logger: logger,
	}
}

func (l *DBLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return &DBLogger{
		logger: l.logger.Desugar().WithOptions(zap.AddCallerSkip(1)).Sugar(),
	}
}

func (l *DBLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Infow(msg, args...)
}

func (l *DBLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Warnw(msg, args...)
}

func (l *DBLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logger.Errorw(msg, args...)
}

func (l *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []interface{}{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.logger.Debugw("trace", fields...)
}
