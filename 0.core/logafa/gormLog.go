package logafa

import (
	"context"
	"fmt"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

type GormLogger struct {
	logLevel gormLogger.LogLevel
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		logLevel: gormLogger.Silent,
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.logLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel != gormLogger.Info {
		return // 非 Debug() 不打印
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		Error("SQL ERR: %v | %v | rows=%d", err, sql, rows)
		return
	}

	fmt.Printf("SQL: %v | rows=%d | time=%v", sql, rows, elapsed)
}
