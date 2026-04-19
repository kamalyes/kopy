/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\errors\codes.go
 * @Description: 错误代码定义 - 基于 go-toolbox/errorx 的类型化错误代码
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errors

import (
	"github.com/kamalyes/go-toolbox/pkg/errorx"
)

// ==================== 模板源错误类型（3000-3004） ====================

const (
	// ErrTypeSourceResolveFailed 解析模板源路径失败
	ErrTypeSourceResolveFailed errorx.ErrorType = 3000 + iota
	// ErrTypeSourceAccessFailed 访问模板源路径失败
	ErrTypeSourceAccessFailed
	// ErrTypeSourceNotDirectory 模板源路径不是目录
	ErrTypeSourceNotDirectory
	// ErrTypeSourceGitCloneFailed Git 克隆仓库失败
	ErrTypeSourceGitCloneFailed
	// ErrTypeSourceTempDirFailed 创建临时目录失败
	ErrTypeSourceTempDirFailed
)

// ==================== 配置错误类型（3100-3107） ====================

const (
	// ErrTypeConfigReadFailed 读取配置文件失败
	ErrTypeConfigReadFailed errorx.ErrorType = 3100 + iota
	// ErrTypeConfigParseFailed 解析配置文件失败
	ErrTypeConfigParseFailed
	// ErrTypeConfigNameRequired 模板名称为必填项
	ErrTypeConfigNameRequired
	// ErrTypeConfigQuestionNameRequired 问题名称为必填项
	ErrTypeConfigQuestionNameRequired
	// ErrTypeConfigQuestionDuplicateName 问题名称重复
	ErrTypeConfigQuestionDuplicateName
	// ErrTypeConfigSelectRequiresChoices 选择类型问题需要提供选项
	ErrTypeConfigSelectRequiresChoices
	// ErrTypeConfigInvalid 配置校验不通过
	ErrTypeConfigInvalid
)

// ==================== 模板引擎错误类型（3200-3202） ====================

const (
	// ErrTypeTemplateParseFailed 解析模板失败
	ErrTypeTemplateParseFailed errorx.ErrorType = 3200 + iota
	// ErrTypeTemplateExecuteFailed 执行模板渲染失败
	ErrTypeTemplateExecuteFailed
	// ErrTypeTemplateRenderFailed 渲染模板失败
	ErrTypeTemplateRenderFailed
)

// ==================== 文件写入错误类型（3300-3303） ====================

const (
	// ErrTypeFileWriteFailed 写入文件失败
	ErrTypeFileWriteFailed errorx.ErrorType = 3300 + iota
	// ErrTypeFileReadFailed 读取文件失败
	ErrTypeFileReadFailed
	// ErrTypeDirCreateFailed 创建目录失败
	ErrTypeDirCreateFailed
	// ErrTypePathResolveFailed 解析路径失败
	ErrTypePathResolveFailed
)

// ==================== 遍历错误类型（3400-3401） ====================

const (
	// ErrTypeWalkRelativePathFailed 获取相对路径失败
	ErrTypeWalkRelativePathFailed errorx.ErrorType = 3400 + iota
	// ErrTypeWalkFailed 遍历模板目录失败
	ErrTypeWalkFailed
)

// init 注册所有错误类型及其默认消息模板
func init() {
	errorx.RegisterError(ErrTypeSourceResolveFailed, "解析模板源路径失败: %s")
	errorx.RegisterError(ErrTypeSourceAccessFailed, "访问模板源路径失败: %s")
	errorx.RegisterError(ErrTypeSourceNotDirectory, "模板源路径不是目录: %s")
	errorx.RegisterError(ErrTypeSourceGitCloneFailed, "Git 克隆仓库失败: %s")
	errorx.RegisterError(ErrTypeSourceTempDirFailed, "创建临时目录失败: %s")

	errorx.RegisterError(ErrTypeConfigReadFailed, "读取配置文件失败: %s")
	errorx.RegisterError(ErrTypeConfigParseFailed, "解析配置文件失败: %s")
	errorx.RegisterError(ErrTypeConfigNameRequired, "模板名称为必填项: %s")
	errorx.RegisterError(ErrTypeConfigQuestionNameRequired, "问题[%s]: 名称不能为空")
	errorx.RegisterError(ErrTypeConfigQuestionDuplicateName, "问题[%s]: 名称重复 %q")
	errorx.RegisterError(ErrTypeConfigSelectRequiresChoices, "问题[%s]: 选择类型需要提供选项")
	errorx.RegisterError(ErrTypeConfigInvalid, "配置校验不通过: %s")

	errorx.RegisterError(ErrTypeTemplateParseFailed, "解析模板失败: %s")
	errorx.RegisterError(ErrTypeTemplateExecuteFailed, "执行模板渲染失败: %s")
	errorx.RegisterError(ErrTypeTemplateRenderFailed, "渲染模板失败: %s")

	errorx.RegisterError(ErrTypeFileWriteFailed, "写入文件失败: %s")
	errorx.RegisterError(ErrTypeFileReadFailed, "读取文件失败: %s")
	errorx.RegisterError(ErrTypeDirCreateFailed, "创建目录失败: %s")
	errorx.RegisterError(ErrTypePathResolveFailed, "解析路径失败: %s")

	errorx.RegisterError(ErrTypeWalkRelativePathFailed, "获取相对路径失败: %s")
	errorx.RegisterError(ErrTypeWalkFailed, "遍历模板目录失败: %s")
}
