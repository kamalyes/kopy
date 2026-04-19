/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\bootstrap\helpers.go
 * @Description: 命令辅助函数 - 参数解析、冲突策略、变量拆分
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bootstrap

import (
	"os"

	kopylog "github.com/kamalyes/kopy/log"
	"github.com/spf13/cobra"
)

func getOutputDir(cmd *cobra.Command) string {
	dir, _ := cmd.Flags().GetString("output")
	return dir
}

func getVarOverrides(cmd *cobra.Command) map[string]string {
	rawVars, _ := cmd.Flags().GetStringSlice("var")
	overrides := make(map[string]string)
	for _, v := range rawVars {
		parts := splitVar(v)
		if len(parts) == 2 {
			overrides[parts[0]] = parts[1]
		}
	}
	return overrides
}

func getConflictStrategy(cmd *cobra.Command) string {
	overwrite, _ := cmd.Flags().GetBool("overwrite")
	skip, _ := cmd.Flags().GetBool("skip")
	if overwrite {
		return "overwrite"
	}
	if skip {
		return "skip"
	}
	return "ask"
}

func getDryRun(cmd *cobra.Command) bool {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	return dryRun
}

func splitVar(s string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return nil
}

func exitWithError(format string, args ...interface{}) {
	kopylog.Error(format, args...)
	os.Exit(1)
}
