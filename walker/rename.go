/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\walker\rename.go
 * @Description: 文件重命名规则 - 对路径应用重命名规则
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package walker

import (
	"strings"

	"github.com/kamalyes/kopy/config"
)

// applyRenameRules 对路径应用重命名规则
func applyRenameRules(path string, rules []config.RenameRule) string {
	for _, rule := range rules {
		path = strings.ReplaceAll(path, rule.Pattern, rule.Value)
	}
	return path
}
