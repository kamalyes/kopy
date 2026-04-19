/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\config\helper.go
 * @Description: 模板配置辅助函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package config

import "github.com/kamalyes/go-toolbox/pkg/mathx"

// GetDefaultOrDefault 安全获取字符串默认值
func GetDefaultOrDefault(val, def string) string {
	return mathx.IfEmpty(val, def)
}
