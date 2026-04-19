/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\engine\render.go
 * @Description: 模板渲染方法 - 字符串、内容、路径渲染
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import (
	"bytes"
	"text/template"

	kopyerrors "github.com/kamalyes/kopy/errors"
)

func (e *Engine) RenderString(input string) (string, error) {
	tmpl, err := template.New("inline").
		Delims(e.cfg.LeftDelimiter(), e.cfg.RightDelimiter()).
		Funcs(e.funcMap).
		Parse(input)
	if err != nil {
		return "", kopyerrors.NewTemplateParseFailedError(err.Error())
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, e.vars); err != nil {
		return "", kopyerrors.NewTemplateExecuteFailedError(err.Error())
	}

	return buf.String(), nil
}

func (e *Engine) RenderContent(content []byte) ([]byte, error) {
	result, err := e.RenderString(string(content))
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

func (e *Engine) RenderPath(path string) (string, error) {
	return e.RenderString(path)
}
