/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\config\validate.go
 * @Description: 模板配置校验 - 校验模板配置合法性
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package config

import (
	"strconv"

	kopyerrors "github.com/kamalyes/kopy/errors"
)

func (c *TemplateConfig) Validate() error {
	if c.Name == "" {
		return kopyerrors.NewConfigNameRequiredError("config")
	}

	names := make(map[string]bool)
	for i, q := range c.Questions {
		idx := strconv.Itoa(i)
		if q.Name == "" {
			return kopyerrors.NewConfigQuestionNameRequiredError(idx)
		}
		if names[q.Name] {
			return kopyerrors.NewConfigQuestionDuplicateNameError(idx, q.Name)
		}
		names[q.Name] = true

		if q.Type == QuestionTypeSelect && len(q.Choices) == 0 {
			return kopyerrors.NewConfigSelectRequiresChoicesError(idx)
		}
	}

	return nil
}
