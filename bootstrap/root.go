/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\bootstrap\root.go
 * @Description: 根命令定义 - 全局参数与命令入口
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bootstrap

import (
	"github.com/spf13/cobra"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "kopy",
	Short: "通用项目模板生成工具 | Universal project template generator",
	Long: `kopy - 通用项目模板生成工具

从本地目录或 Git 仓库加载模板，交互式收集变量，渲染并生成新项目。
支持自定义定界符、忽略规则、命名转换函数（snake/camel/pascal/kebab）等。

示例:
  kopy new ./my-template
  kopy new @kamalyes/gateway-template-service
  kopy new @kamalyes/gateway-template-service@v1.0
  kopy new ./my-template --output ./my-project --overwrite
  kopy new ./my-template --var ServiceName=Payment --var ModuleName=payment`,
}

// Execute 执行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", ".", "输出目录 | Output directory")
	rootCmd.PersistentFlags().StringSliceP("var", "v", []string{}, "模板变量覆盖 (key=value) | Template variable overrides")
	rootCmd.PersistentFlags().Bool("overwrite", false, "覆盖已存在的文件 | Overwrite existing files")
	rootCmd.PersistentFlags().Bool("dry-run", false, "试运行模式 | Dry run mode")
	rootCmd.PersistentFlags().Bool("skip", false, "跳过已存在的文件 | Skip existing files")
}
