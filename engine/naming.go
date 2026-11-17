/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\engine\naming.go
 * @Description: 命名转换 - snake/camel/pascal/kebab，纯标准库实现
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// toSnakeCase 将字符串转换为蛇形命名法（例如：user_name）
// 支持输入格式：camelCase, PascalCase, kebab-case, 普通字符串（带空格）
func toSnakeCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var result strings.Builder
	result.Grow(len(s))

	var prev rune
	lastUnderscore := false

	for i, r := range s {
		// 将连字符和空格都转换为下划线，避免连续下划线
		if r == '-' || r == ' ' {
			if result.Len() > 0 && !lastUnderscore {
				result.WriteByte('_')
				lastUnderscore = true
			}
			prev = r
			continue
		}

		if unicode.IsUpper(r) {
			// 非首字符，且前一个字符不是大写/下划线/空格/连字符时插入下划线
			if i > 0 && !unicode.IsUpper(prev) && prev != '_' && prev != ' ' && prev != '-' {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
			lastUnderscore = false
		} else {
			result.WriteRune(r)
			lastUnderscore = r == '_'
		}
		prev = r
	}

	return result.String()
}

// toPascalCase 将字符串转换为帕斯卡命名法（首字母大写的驼峰，例如：UserName）
// 支持输入格式：camelCase, snake_case, kebab-case, 普通字符串（带空格）
func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	// 包含分隔符时先统一为下划线再分段处理
	if strings.ContainsAny(s, "_- ") {
		normalized := strings.NewReplacer("-", "_", " ", "_").Replace(s)
		parts := strings.Split(normalized, "_")
		var builder strings.Builder
		for _, part := range parts {
			if part == "" {
				continue
			}
			r, size := utf8.DecodeRuneInString(part)
			builder.WriteRune(unicode.ToUpper(r))
			builder.WriteString(strings.ToLower(part[size:]))
		}
		return builder.String()
	}

	// 已经是驼峰式时只需首字母大写
	r, size := utf8.DecodeRuneInString(s)
	if unicode.IsLower(r) {
		return string(unicode.ToUpper(r)) + s[size:]
	}
	return s
}

// toCamelCase 将字符串转换为驼峰命名法（首字母小写，例如：userName）
// 支持输入格式：PascalCase, snake_case, kebab-case
func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if pascal == "" {
		return pascal
	}
	r, size := utf8.DecodeRuneInString(pascal)
	return string(unicode.ToLower(r)) + pascal[size:]
}

// toKebabCase 将字符串转换为短横线命名法（例如：user-name）
// 支持输入格式：camelCase, PascalCase, snake_case
func toKebabCase(s string) string {
	return strings.ReplaceAll(toSnakeCase(s), "_", "-")
}

// titleCase 将每个单词的首字母转换为标题格式
// 与 strings.Title 的单词边界规则保持一致，避免使用已弃用的 API
func titleCase(s string) string {
	prev := ' '
	return strings.Map(func(r rune) rune {
		if isTitleSeparator(prev) {
			prev = r
			return unicode.ToTitle(r)
		}
		prev = r
		return r
	}, s)
}

func isTitleSeparator(r rune) bool {
	if r <= 0x7F {
		switch {
		case '0' <= r && r <= '9':
			return false
		case 'a' <= r && r <= 'z':
			return false
		case 'A' <= r && r <= 'Z':
			return false
		case r == '_':
			return false
		}
		return true
	}
	if unicode.IsLetter(r) || unicode.IsDigit(r) {
		return false
	}
	return unicode.IsSpace(r)
}
