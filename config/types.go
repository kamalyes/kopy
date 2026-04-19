/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\config\types.go
 * @Description: 模板配置类型定义 - 常量、枚举、结构体
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package config

// ConfigFileName 模板配置文件名
const ConfigFileName = "kopy.yaml"

// ===================== 问题类型定义 =====================

// QuestionType 问题类型枚举
type QuestionType string

const (
	// QuestionTypeString 文本输入
	QuestionTypeString QuestionType = "string"
	// QuestionTypeBool 布尔选择
	QuestionTypeBool QuestionType = "bool"
	// QuestionTypeInt 整数输入
	QuestionTypeInt QuestionType = "int"
	// QuestionTypeSelect 下拉选择
	QuestionTypeSelect QuestionType = "select"
)

// ===================== 配置结构体 =====================

// Question 交互式问题定义
type Question struct {
	// Name 变量名（用于模板引用，如 [[.ServiceName]]）
	Name string `yaml:"name"`
	// Prompt 提示文本
	Prompt string `yaml:"prompt"`
	// Type 问题类型
	Type QuestionType `yaml:"type"`
	// Default 默认值（支持模板语法）
	Default string `yaml:"default,omitempty"`
	// Required 是否必填
	Required bool `yaml:"required,omitempty"`
	// Choices 选项列表（仅 select 类型）
	Choices []string `yaml:"choices,omitempty"`
	// Validate 校验规则名称
	Validate string `yaml:"validate,omitempty"`
	// Transform 值转换函数名称
	Transform string `yaml:"transform,omitempty"`
}

// RenameRule 文件/目录重命名规则
type RenameRule struct {
	// Pattern 匹配模式
	Pattern string `yaml:"pattern"`
	// Value 替换值
	Value string `yaml:"value"`
}

// Hook 生命周期钩子命令
type Hook struct {
	// Before 渲染前执行的命令列表
	Before []string `yaml:"before,omitempty"`
	// After 渲染后执行的命令列表
	After []string `yaml:"after,omitempty"`
}

// TemplateConfig 模板配置完整定义
type TemplateConfig struct {
	// Name 模板名称
	Name string `yaml:"name"`
	// Description 模板描述
	Description string `yaml:"description"`
	// Version 模板版本
	Version string `yaml:"version"`
	// Author 模板作者
	Author string `yaml:"author,omitempty"`
	// Tags 模板标签
	Tags []string `yaml:"tags,omitempty"`

	// Questions 交互式问题列表
	Questions []Question `yaml:"questions,omitempty"`
	// Ignore 忽略的文件模式列表
	Ignore []string `yaml:"ignore,omitempty"`
	// Rename 文件重命名规则
	Rename []RenameRule `yaml:"rename,omitempty"`
	// Hooks 生命周期钩子
	Hooks Hook `yaml:"hooks,omitempty"`

	// Delimiters 模板定界符（默认 [[ ]]）
	Delimiters []string `yaml:"delimiters,omitempty"`
}

// LeftDelimiter 获取左定界符
func (c *TemplateConfig) LeftDelimiter() string {
	if len(c.Delimiters) >= 2 {
		return c.Delimiters[0]
	}
	return "[["
}

// RightDelimiter 获取右定界符
func (c *TemplateConfig) RightDelimiter() string {
	if len(c.Delimiters) >= 2 {
		return c.Delimiters[1]
	}
	return "]]"
}
