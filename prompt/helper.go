/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\prompt\helper.go
 * @Description: 提示器辅助方法 - 默认值渲染
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package prompt

import (
	"fmt"
	"strings"
)

// renderDefault 渲染默认值中的模板语法
//
// 失败或渲染出 <no value> 时会给出告警而不是静默放过：默认值引用了尚未提问的变量
// （前向引用）是模板作者的笔误，而一旦静默，未渲染的原文或 <no value> 会被当成答案
// 写进每一个生成文件——曾经导致生成项目的 package.json 里出现 "name": "<no value>"。
func (p *Prompter) renderDefault(name, defaultVal string) string {
	if defaultVal == "" {
		return ""
	}
	rendered, err := p.eng.RenderString(defaultVal)
	if err != nil {
		fmt.Printf("  ⚠ 问题 %s 的默认值渲染失败，按原文使用：%v\n", name, err)
		fmt.Printf("    默认值只能引用排在它前面的变量，且 default 不支持前向引用\n")
		return defaultVal
	}
	if strings.Contains(rendered, "<no value>") {
		fmt.Printf("  ⚠ 问题 %s 的默认值渲染出 <no value>：%q\n", name, rendered)
		fmt.Printf("    它引用了尚未提问的变量，请把被引用的问题排到前面，或改成字面量\n")
	}
	return rendered
}
