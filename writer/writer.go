/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\writer\writer.go
 * @Description: 文件写入器 - 结构体定义、Option 模式、冲突策略
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package writer

import (
	"github.com/kamalyes/kopy/engine"
)

// ===================== 冲突策略定义 =====================

// ConflictStrategy 文件冲突处理策略
type ConflictStrategy string

const (
	// ConflictSkip 跳过已存在的文件
	ConflictSkip ConflictStrategy = "skip"
	// ConflictOverwrite 强制覆盖已存在的文件
	ConflictOverwrite ConflictStrategy = "overwrite"
	// ConflictAsk 交互式询问是否覆盖
	ConflictAsk ConflictStrategy = "ask"
)

// ===================== 写入器定义 =====================

// Writer 文件写入器
type Writer struct {
	outputDir   string
	engine      *engine.Engine
	strategy    ConflictStrategy
	dryRun      bool
	created     int
	skipped     int
	overwritten int
}

// Option 写入器配置选项
type Option func(*Writer)

// WithStrategy 设置冲突处理策略
func WithStrategy(s ConflictStrategy) Option {
	return func(w *Writer) { w.strategy = s }
}

// WithDryRun 设置试运行模式
func WithDryRun(d bool) Option {
	return func(w *Writer) { w.dryRun = d }
}

// New 创建文件写入器
func New(outputDir string, eng *engine.Engine, opts ...Option) *Writer {
	w := &Writer{
		outputDir: outputDir,
		engine:    eng,
		strategy:  ConflictAsk,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Stats 返回写入统计信息
func (w *Writer) Stats() (created, skipped, overwritten int) {
	return w.created, w.skipped, w.overwritten
}
