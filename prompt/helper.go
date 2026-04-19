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

// renderDefault 渲染默认值中的模板语法
func (p *Prompter) renderDefault(defaultVal string) string {
	if defaultVal == "" {
		return ""
	}
	rendered, err := p.eng.RenderString(defaultVal)
	if err != nil {
		return defaultVal
	}
	return rendered
}
