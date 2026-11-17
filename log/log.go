/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\log\log.go
 * @Description: 统一日志管理 - 基于标准库 log 的轻量封装，提供全局日志实例
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Logger 日志接口，仅保留 kopy 实际使用的方法
type Logger interface {
	Debug(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(format string, args ...interface{})
}

const (
	colorReset = "\033[0m"
	colorDebug = "\033[36m"
	colorInfo  = "\033[32m"
	colorWarn  = "\033[33m"
	colorError = "\033[31m"
	colorFatal = "\033[35m"
)

// stdLogger 基于标准库 log 的实现
type stdLogger struct {
	logger *log.Logger
	prefix string
}

var globalLogger Logger = newLogger(os.Stdout, "[KOPY]")

func newLogger(out io.Writer, prefix string) *stdLogger {
	return &stdLogger{
		logger: log.New(out, "", log.LstdFlags),
		prefix: prefix,
	}
}

// WithPrefix 返回一个带指定前缀的新日志实例
func (l *stdLogger) WithPrefix(prefix string) Logger {
	return newLogger(l.logger.Writer(), prefix)
}

func (l *stdLogger) output(color, icon, level, format string, args ...interface{}) {
	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}
	l.logger.Printf("%s %s%s [%s]%s %s\n", l.prefix, color, icon, level, colorReset, msg)
}

func (l *stdLogger) Debug(format string, args ...interface{}) {
	l.output(colorDebug, "🐛", "DEBUG", format, args...)
}

func (l *stdLogger) Info(format string, args ...interface{}) {
	l.output(colorInfo, "ℹ️", "INFO", format, args...)
}

func (l *stdLogger) Warn(format string, args ...interface{}) {
	l.output(colorWarn, "⚠️", "WARN", format, args...)
}

func (l *stdLogger) Error(format string, args ...interface{}) {
	l.output(colorError, "❌", "ERROR", format, args...)
}

func (l *stdLogger) Fatal(format string, args ...interface{}) {
	l.output(colorFatal, "💀", "FATAL", format, args...)
}

// L 获取全局日志实例
func L() Logger {
	return getLogger()
}

var loggerMu sync.RWMutex

// SetLogger 设置自定义日志实例
func SetLogger(l Logger) {
	if l == nil {
		return
	}
	loggerMu.Lock()
	globalLogger = l
	loggerMu.Unlock()
}

// SetOutput 设置全局日志输出目标（测试或嵌入场景使用）
func SetOutput(w io.Writer) {
	if w == nil {
		return
	}
	loggerMu.Lock()
	globalLogger = newLogger(w, "[KOPY]")
	loggerMu.Unlock()
}

func getLogger() Logger {
	loggerMu.RLock()
	l := globalLogger
	loggerMu.RUnlock()
	return l
}

// Debug 记录调试日志
func Debug(format string, args ...interface{}) {
	getLogger().Debug(format, args...)
}

// Info 记录信息日志
func Info(format string, args ...interface{}) {
	getLogger().Info(format, args...)
}

// Warn 记录警告日志
func Warn(format string, args ...interface{}) {
	getLogger().Warn(format, args...)
}

// Error 记录错误日志
func Error(format string, args ...interface{}) {
	getLogger().Error(format, args...)
}

// Fatal 记录致命错误日志
func Fatal(format string, args ...interface{}) {
	getLogger().Fatal(format, args...)
}
