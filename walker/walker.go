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
	"strings"

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
//
// 忽略判定用的是「应用 rename 并渲染之后」的路径，忽略模式本身也先渲染一遍。
// 原先两者都拿未渲染的字符串比对，于是：
//   - rename 的 value 是模板时（value: "[[ if ... ]]__skip__[[ end ]]"），
//     渲染结果即便等于 __skip__，ignore 也命中不了，文件照样生成；
//   - ignore 模式里写模板（"[[ if ne .X \"a\" ]]only-a.txt[[ end ]]"）同样无效。
//
// 这使得「按变量条件性地排除某个目录」在模板里无法表达。
func WalkTemplate(templateDir string, cfg *config.TemplateConfig, w *writer.Writer) (*WalkResult, error) {
	result := &WalkResult{}
	eng := w.Engine()

	// 忽略模式只随变量变化，与遍历到哪个文件无关，故整次遍历只渲染一遍
	ignorePatterns := make([]string, 0, len(cfg.Ignore))
	for _, pattern := range cfg.Ignore {
		rendered, err := eng.RenderString(pattern)
		if err != nil {
			// 渲染不了就按原文比对，保持旧行为，不因一条坏模式中断整次生成
			rendered = pattern
		}
		// 渲染成空串的模式要丢掉：否则空串会与「跳过根目录」之外的空路径意外相等
		if strings.TrimSpace(rendered) == "" {
			continue
		}
		ignorePatterns = append(ignorePatterns, rendered)
	}

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

		// 忽略判定看渲染后的路径；渲染失败则退回原文，与旧行为一致。
		// 这里只影响判定，真正写盘仍由 writer 自己渲染一次，故不改变写入行为
		ignorePath := relPath
		if rendered, err := eng.RenderString(relPath); err == nil {
			ignorePath = rendered
		}

		// 检查忽略规则
		if shouldIgnore(ignorePath, ignorePatterns) {
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
