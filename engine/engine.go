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

// SetVar 写入单个模板变量
//
// 交互式收集变量时每答完一题就写入一个，好让后面问题的 default 能引用前面的答案
// （见 prompt.Collect）
func (e *Engine) SetVar(name string, value interface{}) {
	if e.vars == nil {
		e.vars = make(map[string]interface{})
	}
	e.vars[name] = value
}

// Vars 返回当前模板变量的副本
//
// 返回副本而非内部 map：调用方拿去做种子或打印时改动它不应影响引擎状态
func (e *Engine) Vars() map[string]interface{} {
	out := make(map[string]interface{}, len(e.vars))
	for k, v := range e.vars {
		out[k] = v
	}
	return out
}
