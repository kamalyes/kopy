/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\config\loader.go
 * @Description: 模板配置加载 - 从目录读取并解析 kopy.yaml
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package config

import (
	"os"
	"path/filepath"

	kopyerrors "github.com/kamalyes/kopy/errors"
	"gopkg.in/yaml.v3"
)

func LoadFromDir(dir string) (*TemplateConfig, error) {
	configPath := filepath.Join(dir, ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, kopyerrors.NewConfigReadFailedError(configPath)
	}

	cfg := &TemplateConfig{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, kopyerrors.NewConfigParseFailedError(configPath)
	}

	if cfg.Name == "" {
		return nil, kopyerrors.NewConfigNameRequiredError(configPath)
	}

	if len(cfg.Delimiters) == 0 {
		cfg.Delimiters = []string{"[[", "]]"}
	}

	if len(cfg.Ignore) == 0 {
		cfg.Ignore = []string{ConfigFileName, ".git/**"}
	}

	return cfg, nil
}
