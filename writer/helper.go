/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\writer\helper.go
 * @Description: 写入器辅助函数 - 冲突询问、可执行文件检测
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package writer

import (
	"fmt"
	"strings"
)

// askOverwrite 交互式询问是否覆盖
func askOverwrite(path string) bool {
	fmt.Printf("  ⚠  %s already exists. Overwrite? [y/N]: ", path)
	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

// isExecutable 判断内容是否为可执行文件（shebang 检测）
func isExecutable(content []byte) bool {
	if len(content) < 2 {
		return false
	}
	return content[0] == '#' && content[1] == '!'
}
