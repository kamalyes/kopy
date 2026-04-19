/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\prompt\prompt.go
 * @Description: 交互式提示器 - 结构体定义与变量收集
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package prompt

import (
	"bufio"
	"fmt"
	"os"

	"github.com/kamalyes/kopy/config"
	"github.com/kamalyes/kopy/engine"
)

// Prompter 交互式提示器
type Prompter struct {
	reader *bufio.Reader
	eng    *engine.Engine
}

// New 创建交互式提示器
func New(eng *engine.Engine) *Prompter {
	return &Prompter{
		reader: bufio.NewReader(os.Stdin),
		eng:    eng,
	}
}

// Collect 收集所有问题的答案，overrides 中的值将跳过交互
//
// 每答完一题就把结果写回引擎，这样后面问题的 default 才能引用前面的答案，
// 例如 ProjectName 的 default 写 "[[ .OrgPrefix ]]-share-proto"。
// 原先是全部收集完才 SetVars，渲染 default 时引擎里一个变量都没有：
// 引用别的变量得到 <no value>，带过滤器的（如 | lower）还会在 Execute 阶段
// 因 nil 入参报错，被 renderDefault 的兜底分支原文返回，于是未渲染的模板串
// 直接变成答案再写进每个生成文件。
//
// 种子取引擎当前已有的变量（调用方可预置 Year 等内置值）。
func (p *Prompter) Collect(questions []config.Question, overrides map[string]string) (map[string]interface{}, error) {
	vars := p.eng.Vars()

	for _, q := range questions {
		// 命令行覆盖优先
		if val, ok := overrides[q.Name]; ok {
			vars[q.Name] = val
			p.eng.SetVar(q.Name, val)
			continue
		}

		val, err := p.ask(q)
		if err != nil {
			return nil, fmt.Errorf("question %q: %w", q.Name, err)
		}
		vars[q.Name] = val
		p.eng.SetVar(q.Name, val)
	}

	return vars, nil
}

// ask 根据问题类型分发到对应的输入方法
func (p *Prompter) ask(q config.Question) (string, error) {
	defaultVal := p.renderDefault(q.Name, q.Default)

	for {
		display := q.Prompt
		if display == "" {
			display = q.Name
		}

		var val string
		var err error

		switch q.Type {
		case config.QuestionTypeSelect:
			val, err = p.askSelect(display, q.Choices, defaultVal)
		case config.QuestionTypeBool:
			val, err = p.askBool(display, defaultVal)
		case config.QuestionTypeInt:
			val, err = p.askInt(display, defaultVal)
		default:
			val, err = p.askString(display, defaultVal)
		}

		if err != nil {
			return "", err
		}

		// 必填校验
		if q.Required && val == "" {
			fmt.Println("  ⚠ This field is required")
			continue
		}

		return val, nil
	}
}
