/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\log\log.go
 * @Description: 统一日志管理 - 基于 go-logger 封装，提供全局日志实例
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package log

import (
	"github.com/kamalyes/go-logger"
)

var globalLogger logger.ILogger

func init() {
	globalLogger = logger.New().WithPrefix("[KOPY]")
}

// L 获取全局日志实例
func L() logger.ILogger {
	return globalLogger
}

// SetLogger 设置自定义日志实例
func SetLogger(l logger.ILogger) {
	if l != nil {
		globalLogger = l
	}
}

// Debug 记录调试日志
func Debug(format string, args ...interface{}) {
	globalLogger.Debug(format, args...)
}

// Info 记录信息日志
func Info(format string, args ...interface{}) {
	globalLogger.Info(format, args...)
}

// Warn 记录警告日志
func Warn(format string, args ...interface{}) {
	globalLogger.Warn(format, args...)
}

// Error 记录错误日志
func Error(format string, args ...interface{}) {
	globalLogger.Error(format, args...)
}

// Fatal 记录致命错误日志
func Fatal(format string, args ...interface{}) {
	globalLogger.Fatal(format, args...)
}
