/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\writer\write.go
 * @Description: 文件写入方法 - 渲染模板并写入目标目录，支持冲突检测
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package writer

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFile 渲染并写入单个文件
func (w *Writer) WriteFile(relPath string, content []byte) error {
	// 渲染文件路径
	renderedPath, err := w.engine.RenderPath(relPath)
	if err != nil {
		return fmt.Errorf("render path %s: %w", relPath, err)
	}

	// 渲染文件内容
	renderedContent, err := w.engine.RenderContent(content)
	if err != nil {
		return fmt.Errorf("render content %s: %w", relPath, err)
	}

	fullPath := filepath.Join(w.outputDir, renderedPath)

	// 试运行模式
	if w.dryRun {
		fmt.Printf("  [DRY] %s\n", renderedPath)
		w.created++
		return nil
	}

	// 冲突检测
	if _, err := os.Stat(fullPath); err == nil {
		switch w.strategy {
		case ConflictSkip:
			fmt.Printf("  ⏭  %s (exists, skipped)\n", renderedPath)
			w.skipped++
			return nil
		case ConflictOverwrite:
			fmt.Printf("  🔄 %s (overwritten)\n", renderedPath)
			w.overwritten++
		case ConflictAsk:
			if !askOverwrite(renderedPath) {
				fmt.Printf("  ⏭  %s (exists, skipped)\n", renderedPath)
				w.skipped++
				return nil
			}
			fmt.Printf("  🔄 %s (overwritten)\n", renderedPath)
			w.overwritten++
		}
	} else {
		fmt.Printf("  ✅ %s\n", renderedPath)
		w.created++
	}

	// 确保目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dir %s: %w", dir, err)
	}

	// 根据内容判断文件权限
	perm := os.FileMode(0644)
	if isExecutable(content) {
		perm = 0755
	}

	if err := os.WriteFile(fullPath, renderedContent, perm); err != nil {
		return fmt.Errorf("write file %s: %w", fullPath, err)
	}

	return nil
}

// CreateDir 创建目录
func (w *Writer) CreateDir(relPath string) error {
	renderedPath, err := w.engine.RenderPath(relPath)
	if err != nil {
		return fmt.Errorf("render path %s: %w", relPath, err)
	}

	fullPath := filepath.Join(w.outputDir, renderedPath)

	if w.dryRun {
		fmt.Printf("  [DRY] %s/\n", renderedPath)
		return nil
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", fullPath, err)
		}
	}
	return nil
}
