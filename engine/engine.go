/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\engine\engine.go
 * @Description: 模板渲染引擎 - 结构体定义与构造
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import (
	"text/template"

	"github.com/kamalyes/kopy/config"
)

// Engine 模板渲染引擎
type Engine struct {
	cfg     *config.TemplateConfig
	vars    map[string]interface{}
	funcMap template.FuncMap
}

// New 创建模板渲染引擎
func New(cfg *config.TemplateConfig) *Engine {
	e := &Engine{
		cfg:  cfg,
		vars: make(map[string]interface{}),
	}
	e.funcMap = e.buildFuncMap()
	return e
}

// SetVars 设置模板变量
func (e *Engine) SetVars(vars map[string]interface{}) {
	e.vars = vars
}
