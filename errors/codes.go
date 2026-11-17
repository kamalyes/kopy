/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\errors\codes.go
 * @Description: 错误代码定义 - 基于标准库的类型化错误代码
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errors

// ==================== 模板源错误类型（3000-3004） ====================

const (
	// ErrTypeSourceResolveFailed 解析模板源路径失败
	ErrTypeSourceResolveFailed ErrorType = 3000 + iota
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
	ErrTypeConfigReadFailed ErrorType = 3100 + iota
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
	ErrTypeTemplateParseFailed ErrorType = 3200 + iota
	// ErrTypeTemplateExecuteFailed 执行模板渲染失败
	ErrTypeTemplateExecuteFailed
	// ErrTypeTemplateRenderFailed 渲染模板失败
	ErrTypeTemplateRenderFailed
)

// ==================== 文件写入错误类型（3300-3303） ====================

const (
	// ErrTypeFileWriteFailed 写入文件失败
	ErrTypeFileWriteFailed ErrorType = 3300 + iota
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
	ErrTypeWalkRelativePathFailed ErrorType = 3400 + iota
	// ErrTypeWalkFailed 遍历模板目录失败
	ErrTypeWalkFailed
)

// init 注册所有错误类型及其默认消息模板
func init() {
	RegisterError(ErrTypeSourceResolveFailed, "解析模板源路径失败: %s")
	RegisterError(ErrTypeSourceAccessFailed, "访问模板源路径失败: %s")
	RegisterError(ErrTypeSourceNotDirectory, "模板源路径不是目录: %s")
	RegisterError(ErrTypeSourceGitCloneFailed, "Git 克隆仓库失败: %s")
	RegisterError(ErrTypeSourceTempDirFailed, "创建临时目录失败: %s")

	RegisterError(ErrTypeConfigReadFailed, "读取配置文件失败: %s")
	RegisterError(ErrTypeConfigParseFailed, "解析配置文件失败: %s")
	RegisterError(ErrTypeConfigNameRequired, "模板名称为必填项: %s")
	RegisterError(ErrTypeConfigQuestionNameRequired, "问题[%s]: 名称不能为空")
	RegisterError(ErrTypeConfigQuestionDuplicateName, "问题[%s]: 名称重复 %q")
	RegisterError(ErrTypeConfigSelectRequiresChoices, "问题[%s]: 选择类型需要提供选项")
	RegisterError(ErrTypeConfigInvalid, "配置校验不通过: %s")

	RegisterError(ErrTypeTemplateParseFailed, "解析模板失败: %s")
	RegisterError(ErrTypeTemplateExecuteFailed, "执行模板渲染失败: %s")
	RegisterError(ErrTypeTemplateRenderFailed, "渲染模板失败: %s")

	RegisterError(ErrTypeFileWriteFailed, "写入文件失败: %s")
	RegisterError(ErrTypeFileReadFailed, "读取文件失败: %s")
	RegisterError(ErrTypeDirCreateFailed, "创建目录失败: %s")
	RegisterError(ErrTypePathResolveFailed, "解析路径失败: %s")

	RegisterError(ErrTypeWalkRelativePathFailed, "获取相对路径失败: %s")
	RegisterError(ErrTypeWalkFailed, "遍历模板目录失败: %s")
}
