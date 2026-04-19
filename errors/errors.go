/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\errors\errors.go
 * @Description: 统一错误构造 - 基于 go-toolbox/errorx 的错误构造函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errors

import (
	"github.com/kamalyes/go-toolbox/pkg/errorx"
)

// ==================== 模板源错误构造函数 ====================

// NewSourceResolveFailedError 创建解析模板源路径失败错误
func NewSourceResolveFailedError(detail string) error {
	return errorx.NewError(ErrTypeSourceResolveFailed, detail)
}

// NewSourceAccessFailedError 创建访问模板源路径失败错误
func NewSourceAccessFailedError(detail string) error {
	return errorx.NewError(ErrTypeSourceAccessFailed, detail)
}

// NewSourceNotDirectoryError 创建模板源路径不是目录错误
func NewSourceNotDirectoryError(detail string) error {
	return errorx.NewError(ErrTypeSourceNotDirectory, detail)
}

// NewSourceGitCloneFailedError 创建 Git 克隆仓库失败错误
func NewSourceGitCloneFailedError(detail string) error {
	return errorx.NewError(ErrTypeSourceGitCloneFailed, detail)
}

// NewSourceTempDirFailedError 创建临时目录失败错误
func NewSourceTempDirFailedError(detail string) error {
	return errorx.NewError(ErrTypeSourceTempDirFailed, detail)
}

// ==================== 配置错误构造函数 ====================

// NewConfigReadFailedError 创建读取配置文件失败错误
func NewConfigReadFailedError(detail string) error {
	return errorx.NewError(ErrTypeConfigReadFailed, detail)
}

// NewConfigParseFailedError 创建解析配置文件失败错误
func NewConfigParseFailedError(detail string) error {
	return errorx.NewError(ErrTypeConfigParseFailed, detail)
}

// NewConfigNameRequiredError 创建模板名称为必填项错误
func NewConfigNameRequiredError(detail string) error {
	return errorx.NewError(ErrTypeConfigNameRequired, detail)
}

// NewConfigQuestionNameRequiredError 创建问题名称为必填项错误
func NewConfigQuestionNameRequiredError(index string) error {
	return errorx.NewError(ErrTypeConfigQuestionNameRequired, index)
}

// NewConfigQuestionDuplicateNameError 创建问题名称重复错误
func NewConfigQuestionDuplicateNameError(index, name string) error {
	return errorx.NewError(ErrTypeConfigQuestionDuplicateName, index, name)
}

// NewConfigSelectRequiresChoicesError 创建选择类型问题需要提供选项错误
func NewConfigSelectRequiresChoicesError(index string) error {
	return errorx.NewError(ErrTypeConfigSelectRequiresChoices, index)
}

// NewConfigInvalidError 创建配置校验不通过错误
func NewConfigInvalidError(detail string) error {
	return errorx.NewError(ErrTypeConfigInvalid, detail)
}

// ==================== 模板引擎错误构造函数 ====================

// NewTemplateParseFailedError 创建解析模板失败错误
func NewTemplateParseFailedError(detail string) error {
	return errorx.NewError(ErrTypeTemplateParseFailed, detail)
}

// NewTemplateExecuteFailedError 创建执行模板渲染失败错误
func NewTemplateExecuteFailedError(detail string) error {
	return errorx.NewError(ErrTypeTemplateExecuteFailed, detail)
}

// NewTemplateRenderFailedError 创建渲染模板失败错误
func NewTemplateRenderFailedError(detail string) error {
	return errorx.NewError(ErrTypeTemplateRenderFailed, detail)
}

// ==================== 文件写入错误构造函数 ====================

// NewFileWriteFailedError 创建写入文件失败错误
func NewFileWriteFailedError(detail string) error {
	return errorx.NewError(ErrTypeFileWriteFailed, detail)
}

// NewFileReadFailedError 创建读取文件失败错误
func NewFileReadFailedError(detail string) error {
	return errorx.NewError(ErrTypeFileReadFailed, detail)
}

// NewDirCreateFailedError 创建创建目录失败错误
func NewDirCreateFailedError(detail string) error {
	return errorx.NewError(ErrTypeDirCreateFailed, detail)
}

// NewPathResolveFailedError 创建解析路径失败错误
func NewPathResolveFailedError(detail string) error {
	return errorx.NewError(ErrTypePathResolveFailed, detail)
}

// ==================== 遍历错误构造函数 ====================

// NewWalkRelativePathFailedError 创建获取相对路径失败错误
func NewWalkRelativePathFailedError(detail string) error {
	return errorx.NewError(ErrTypeWalkRelativePathFailed, detail)
}

// NewWalkFailedError 创建遍历模板目录失败错误
func NewWalkFailedError(detail string) error {
	return errorx.NewError(ErrTypeWalkFailed, detail)
}

// ==================== 通用错误工具函数 ====================

// WrapError 包装错误，添加上下文信息
func WrapError(message string, err error) error {
	return errorx.WrapError(message, err)
}
