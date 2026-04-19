/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\source\local.go
 * @Description: 本地目录模板源 - 解析本地路径为绝对路径
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package source

import (
	"os"
	"path/filepath"

	kopyerrors "github.com/kamalyes/kopy/errors"
)

type localSource struct {
	path string
}

func (s *localSource) Resolve() (string, error) {
	abs, err := filepath.Abs(s.path)
	if err != nil {
		return "", kopyerrors.NewSourceResolveFailedError(s.path)
	}

	info, err := os.Stat(abs)
	if err != nil {
		return "", kopyerrors.NewSourceAccessFailedError(abs)
	}

	if !info.IsDir() {
		return "", kopyerrors.NewSourceNotDirectoryError(abs)
	}

	return abs, nil
}

func (s *localSource) Cleanup() {}
