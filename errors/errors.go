/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\errors\errors.go
 * @Description: 统一错误构造 - 基于标准库 fmt/errors 的错误构造函数
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errors

// ==================== 模板源错误构造函数 ====================

// NewSourceResolveFailedError 创建解析模板源路径失败错误
func NewSourceResolveFailedError(detail string) error {
	return NewError(ErrTypeSourceResolveFailed, detail)
}

// NewSourceAccessFailedError 创建访问模板源路径失败错误
func NewSourceAccessFailedError(detail string) error {
	return NewError(ErrTypeSourceAccessFailed, detail)
}

// NewSourceNotDirectoryError 创建模板源路径不是目录错误
func NewSourceNotDirectoryError(detail string) error {
	return NewError(ErrTypeSourceNotDirectory, detail)
}

// NewSourceGitCloneFailedError 创建 Git 克隆仓库失败错误
func NewSourceGitCloneFailedError(detail string) error {
	return NewError(ErrTypeSourceGitCloneFailed, detail)
}

// NewSourceTempDirFailedError 创建临时目录失败错误
func NewSourceTempDirFailedError(detail string) error {
	return NewError(ErrTypeSourceTempDirFailed, detail)
}

// ==================== 配置错误构造函数 ====================

// NewConfigReadFailedError 创建读取配置文件失败错误
func NewConfigReadFailedError(detail string) error {
	return NewError(ErrTypeConfigReadFailed, detail)
}

// NewConfigParseFailedError 创建解析配置文件失败错误
func NewConfigParseFailedError(detail string) error {
	return NewError(ErrTypeConfigParseFailed, detail)
}

// NewConfigNameRequiredError 创建模板名称为必填项错误
func NewConfigNameRequiredError(detail string) error {
	return NewError(ErrTypeConfigNameRequired, detail)
}

// NewConfigQuestionNameRequiredError 创建问题名称为必填项错误
func NewConfigQuestionNameRequiredError(index string) error {
	return NewError(ErrTypeConfigQuestionNameRequired, index)
}

// NewConfigQuestionDuplicateNameError 创建问题名称重复错误
func NewConfigQuestionDuplicateNameError(index, name string) error {
	return NewError(ErrTypeConfigQuestionDuplicateName, index, name)
}

// NewConfigSelectRequiresChoicesError 创建选择类型问题需要提供选项错误
func NewConfigSelectRequiresChoicesError(index string) error {
	return NewError(ErrTypeConfigSelectRequiresChoices, index)
}

// NewConfigInvalidError 创建配置校验不通过错误
func NewConfigInvalidError(detail string) error {
	return NewError(ErrTypeConfigInvalid, detail)
}

// ==================== 模板引擎错误构造函数 ====================

// NewTemplateParseFailedError 创建解析模板失败错误
func NewTemplateParseFailedError(detail string) error {
	return NewError(ErrTypeTemplateParseFailed, detail)
}

// NewTemplateExecuteFailedError 创建执行模板渲染失败错误
func NewTemplateExecuteFailedError(detail string) error {
	return NewError(ErrTypeTemplateExecuteFailed, detail)
}

// NewTemplateRenderFailedError 创建渲染模板失败错误
func NewTemplateRenderFailedError(detail string) error {
	return NewError(ErrTypeTemplateRenderFailed, detail)
}

// ==================== 文件写入错误构造函数 ====================

// NewFileWriteFailedError 创建写入文件失败错误
func NewFileWriteFailedError(detail string) error {
	return NewError(ErrTypeFileWriteFailed, detail)
}

// NewFileReadFailedError 创建读取文件失败错误
func NewFileReadFailedError(detail string) error {
	return NewError(ErrTypeFileReadFailed, detail)
}

// NewDirCreateFailedError 创建创建目录失败错误
func NewDirCreateFailedError(detail string) error {
	return NewError(ErrTypeDirCreateFailed, detail)
}

// NewPathResolveFailedError 创建解析路径失败错误
func NewPathResolveFailedError(detail string) error {
	return NewError(ErrTypePathResolveFailed, detail)
}

// ==================== 遍历错误构造函数 ====================

// NewWalkRelativePathFailedError 创建获取相对路径失败错误
func NewWalkRelativePathFailedError(detail string) error {
	return NewError(ErrTypeWalkRelativePathFailed, detail)
}

// NewWalkFailedError 创建遍历模板目录失败错误
func NewWalkFailedError(detail string) error {
	return NewError(ErrTypeWalkFailed, detail)
}
