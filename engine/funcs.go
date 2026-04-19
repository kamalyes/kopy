/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\engine\funcs.go
 * @Description: 内置模板函数 - 命名转换、字符串操作、默认值
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import (
	"strings"
	"text/template"
	"time"

	"github.com/kamalyes/go-toolbox/pkg/stringx"
)

// buildFuncMap 构建模板函数映射，复用 go-toolbox/stringx 的命名转换
func (e *Engine) buildFuncMap() template.FuncMap {
	return template.FuncMap{
		// 命名转换（复用 go-toolbox/stringx）
		"snake":  stringx.ToSnakeCase,
		"camel":  stringx.ToCamelCase,
		"pascal": stringx.ToPascalCase,
		"kebab":  stringx.ToKebabCase,

		// 字符串操作
		"lower":   strings.ToLower,
		"upper":   strings.ToUpper,
		"title":   strings.Title,
		"trim":    strings.TrimSpace,
		"prefix":  func(s, p string) string { return p + s },
		"suffix":  func(s, sfx string) string { return s + sfx },
		"replace": func(s, old, new string) string { return strings.ReplaceAll(s, old, new) },

		// 时间函数
		"now": func() time.Time {
			return time.Now()
		},
		"date": func(layout string, value interface{}) string {
			switch v := value.(type) {
			case time.Time:
				return v.Format(layout)
			case *time.Time:
				if v == nil {
					return ""
				}
				return v.Format(layout)
			case string:
				for _, parseLayout := range []string{
					time.RFC3339Nano,
					time.RFC3339,
					time.DateTime,
					time.DateOnly,
				} {
					if t, err := time.Parse(parseLayout, v); err == nil {
						return t.Format(layout)
					}
				}
				return v
			default:
				return ""
			}
		},
		"nowFmt": func(layout string) string {
			return time.Now().Format(layout)
		},

		// 默认值
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
	}
}
