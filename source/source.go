/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\source\source.go
 * @Description: 模板源接口与工厂 - 自动判断本地/Git 源与 @owner/repo 简写
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
//
// 支持：
//   - 本地目录：./my-template、/path/to/template
//   - 完整 Git 地址：https://github.com/user/template、git@github.com:user/template.git
//   - GitHub 简写：@owner/repo、@owner/repo@branch、@owner/repo#branch
func New(raw string) Source {
	if url, branch, ok := parseShorthand(raw); ok {
		if branch != "" {
			return &gitSource{url: url + "#" + branch}
		}
		return &gitSource{url: url}
	}

	if isGitURL(raw) {
		return &gitSource{url: raw}
	}
	return &localSource{path: raw}
}

// parseShorthand 解析 @owner/repo 与 @owner/repo@branch / @owner/repo#branch
func parseShorthand(raw string) (url, branch string, ok bool) {
	if !strings.HasPrefix(raw, "@") || strings.HasPrefix(raw, "@@") {
		return "", "", false
	}

	body := strings.TrimPrefix(raw, "@")
	repo := body
	if i := strings.IndexAny(body, "@#"); i >= 0 {
		repo = body[:i]
		branch = body[i+1:]
	}

	// 只接受 owner/repo 两段，避免把本地目录误判成 GitHub 简写
	parts := strings.Split(repo, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	return "https://github.com/" + repo, branch, true
}

// isGitURL 判断是否为 Git 仓库地址
func isGitURL(s string) bool {
	return strings.HasPrefix(s, "https://") ||
		strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "git@") ||
		strings.HasPrefix(s, "git://") ||
		strings.HasSuffix(s, ".git")
}
