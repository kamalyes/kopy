/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\bootstrap\list.go
 * @Description: list 子命令 - 列出模板中的问题和默认值
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bootstrap

import (
	"github.com/kamalyes/kopy/config"
	kopylog "github.com/kamalyes/kopy/log"
	"github.com/kamalyes/kopy/source"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list <template-source>",
	Short: "列出模板信息 | List template info and questions",
	Args:  cobra.ExactArgs(1),
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	templateSource := args[0]

	src := source.New(templateSource)
	templateDir, err := src.Resolve()
	if err != nil {
		exitWithError("解析模板源失败: %v", err)
	}
	defer src.Cleanup()

	cfg, err := config.LoadFromDir(templateDir)
	if err != nil {
		exitWithError("加载配置失败: %v", err)
	}

	kopylog.Info("模板名称: %s", cfg.Name)
	kopylog.Info("模板描述: %s", cfg.Description)
	kopylog.Info("模板版本: %s", cfg.Version)
	kopylog.Info("定界符: %s %s", cfg.LeftDelimiter(), cfg.RightDelimiter())

	if len(cfg.Questions) > 0 {
		kopylog.Info("问题列表:")
		for i, q := range cfg.Questions {
			kopylog.Info("  %d. %s (%s) [%s]", i+1, q.Prompt, q.Name, q.Type)
			if q.Default != "" {
				kopylog.Info("     默认值: %s", q.Default)
			}
			if len(q.Choices) > 0 {
				kopylog.Info("     选项: %v", q.Choices)
			}
			if q.Required {
				kopylog.Info("     必填: 是")
			}
		}
	}

	if len(cfg.Ignore) > 0 {
		kopylog.Info("忽略规则:")
		for _, pattern := range cfg.Ignore {
			kopylog.Info("  - %s", pattern)
		}
	}
}
