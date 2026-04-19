/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\walker\walker.go
 * @Description: 模板目录遍历器 - 递归遍历模板目录并过滤忽略文件
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package walker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kamalyes/kopy/config"
	"github.com/kamalyes/kopy/writer"
)

// WalkResult 遍历结果统计
type WalkResult struct {
	Files       int
	Directories int
	Skipped     int
}

// WalkTemplate 遍历模板目录，将所有文件通过 writer 写入目标
func WalkTemplate(templateDir string, cfg *config.TemplateConfig, w *writer.Writer) (*WalkResult, error) {
	result := &WalkResult{}

	err := filepath.Walk(templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 使用正斜杠统一路径分隔符
		relPath = filepath.ToSlash(relPath)

		// 应用重命名规则
		relPath = applyRenameRules(relPath, cfg.Rename)

		// 检查忽略规则
		if shouldIgnore(relPath, cfg.Ignore) {
			if info.IsDir() {
				result.Skipped++
				return filepath.SkipDir
			}
			result.Skipped++
			return nil
		}

		// 处理目录
		if info.IsDir() {
			result.Directories++
			return w.CreateDir(relPath)
		}

		// 处理文件
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s: %w", path, err)
		}

		if err := w.WriteFile(relPath, content); err != nil {
			return fmt.Errorf("write file %s: %w", relPath, err)
		}

		result.Files++
		return nil
	})

	return result, err
}
