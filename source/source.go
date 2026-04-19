/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\source\source.go
 * @Description: 模板源接口与工厂 - 自动判断本地/Git源
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package source

import "strings"

// Source 模板源接口
type Source interface {
	// Resolve 解析模板源路径，返回本地目录绝对路径
	Resolve() (string, error)
	// Cleanup 清理临时资源
	Cleanup()
}

// New 根据输入自动判断创建本地源或 Git 源
func New(raw string) Source {
	if isGitURL(raw) {
		return &gitSource{url: raw}
	}
	return &localSource{path: raw}
}

// isGitURL 判断是否为 Git 仓库地址
func isGitURL(s string) bool {
	return strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "git@") ||
		strings.HasPrefix(s, "git://") ||
		strings.HasSuffix(s, ".git")
}
