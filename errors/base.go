/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\errors\base.go
 * @Description: 类型化错误基础实现 - 标准库 errors/fmt 封装
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errors

import (
	stderrors "errors"
	"fmt"
	"sync"
)

// ErrorType 错误类型编号
type ErrorType int

// BaseError 带类型编号的错误
type BaseError struct {
	Msg  string
	Type ErrorType
}

// Error 实现 error 接口
func (e BaseError) Error() string {
	return e.Msg
}

// GetType 获取错误类型
func (e BaseError) GetType() ErrorType {
	return e.Type
}

var (
	mu       sync.RWMutex
	messages = make(map[ErrorType]string)
)

// RegisterError 注册错误类型及其消息模板
func RegisterError(errType ErrorType, msg string) {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := messages[errType]; exists {
		return
	}
	messages[errType] = msg
}

// NewError 按错误类型和参数创建错误
func NewError(errType ErrorType, args ...interface{}) BaseError {
	mu.RLock()
	msg, ok := messages[errType]
	mu.RUnlock()
	if !ok {
		return BaseError{Msg: "unknown error"}
	}
	return BaseError{Msg: fmt.Sprintf(msg, args...), Type: errType}
}

// WrapError 包装错误并添加上下文
func WrapError(message string, err ...error) error {
	if len(err) == 0 || err[0] == nil {
		return stderrors.New(message)
	}
	return fmt.Errorf("%s: %w", message, err[0])
}
