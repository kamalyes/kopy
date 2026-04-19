/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\walker\ignore.go
 * @Description: 忽略规则匹配 - 支持精确/通配符/前后缀/路径段匹配
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package walker

import (
	"path/filepath"
	"strings"
)

// shouldIgnore 检查路径是否匹配忽略规则
func shouldIgnore(path string, patterns []string) bool {
	for _, pattern := range patterns {
		// 精确匹配
		if path == pattern {
			return true
		}

		// ** 通配符匹配
		if strings.Contains(pattern, "**") {
			globPattern := strings.ReplaceAll(pattern, "**", "*")
			matched, _ := filepath.Match(globPattern, path)
			if matched {
				return true
			}
		}

		// 前缀匹配（目录前缀）
		if strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(path, pattern) || strings.HasPrefix(path+"/", pattern) {
				return true
			}
		}

		// 后缀匹配
		if strings.HasPrefix(pattern, "*.") {
			ext := pattern[1:]
			if strings.HasSuffix(path, ext) {
				return true
			}
		}

		// glob 匹配
		matched, _ := filepath.Match(pattern, path)
		if matched {
			return true
		}

		// 路径中任意段匹配
		parts := strings.Split(path, "/")
		for _, part := range parts {
			matched, _ := filepath.Match(pattern, part)
			if matched {
				return true
			}
		}
	}

	return false
}
