/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\bootstrap\new.go
 * @Description: new 子命令 - 从模板生成新项目
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bootstrap

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kamalyes/kopy/config"
	"github.com/kamalyes/kopy/engine"
	kopylog "github.com/kamalyes/kopy/log"
	"github.com/kamalyes/kopy/prompt"
	"github.com/kamalyes/kopy/source"
	"github.com/kamalyes/kopy/walker"
	"github.com/kamalyes/kopy/writer"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <template-source>",
	Short: "从模板生成新项目 | Generate a new project from template",
	Long: `从本地目录或 Git 仓库加载模板，交互式收集变量，生成新项目。

模板源支持:
  - 本地目录路径: ./my-template, /path/to/template
  - Git 仓库地址: https://github.com/user/template
  - Git 仓库+分支: https://github.com/user/template#v1.0`,
	Args: cobra.ExactArgs(1),
	Run:  runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) {
	templateSource := args[0]
	outputDir := getOutputDir(cmd)
	varOverrides := getVarOverrides(cmd)
	strategy := getConflictStrategy(cmd)
	dryRun := getDryRun(cmd)

	kopylog.Info("🚀 kopy - 通用项目模板生成工具")

	kopylog.Info("📦 解析模板源...")
	src := source.New(templateSource)
	templateDir, err := src.Resolve()
	if err != nil {
		exitWithError("解析模板源失败: %v", err)
	}
	defer src.Cleanup()

	kopylog.Info("📁 模板路径: %s", templateDir)

	kopylog.Info("📋 加载模板配置...")
	cfg, err := config.LoadFromDir(templateDir)
	if err != nil {
		exitWithError("加载配置失败: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		exitWithError("配置校验失败: %v", err)
	}

	kopylog.Info("📝 模板名称: %s", cfg.Name)
	kopylog.Info("📝 模板描述: %s", cfg.Description)
	kopylog.Info("📝 模板版本: %s", cfg.Version)

	eng := engine.New(cfg)

	kopylog.Info("❓ 收集模板变量...")
	p := prompt.New(eng)
	vars, err := p.Collect(cfg.Questions, varOverrides)
	if err != nil {
		exitWithError("收集变量失败: %v", err)
	}

	vars["Year"] = time.Now().Year()
	eng.SetVars(vars)

	kopylog.Info("📝 模板变量:")
	for k, v := range vars {
		kopylog.Info("  %s = %v", k, v)
	}

	renderedOutputDir, err := eng.RenderString(outputDir)
	if err != nil {
		renderedOutputDir = outputDir
	}

	if !filepath.IsAbs(renderedOutputDir) {
		cwd, _ := os.Getwd()
		renderedOutputDir = filepath.Join(cwd, renderedOutputDir)
	}

	kopylog.Info("📂 输出目录: %s", renderedOutputDir)

	w := writer.New(renderedOutputDir, eng,
		writer.WithStrategy(writer.ConflictStrategy(strategy)),
		writer.WithDryRun(dryRun),
	)

	kopylog.Info("🔧 生成项目文件...")
	result, err := walker.WalkTemplate(templateDir, cfg, w)
	if err != nil {
		exitWithError("生成文件失败: %v", err)
	}

	created, skipped, overwritten := w.Stats()
	kopylog.Info("✅ 项目生成完成!")
	kopylog.Info("📄 文件: %d 创建, %d 跳过, %d 覆盖", created, skipped, overwritten)
	kopylog.Info("📁 目录: %d 创建", result.Directories)
	kopylog.Info("⏭  忽略: %d", result.Skipped)
	kopylog.Info("📂 项目路径: %s", renderedOutputDir)
	kopylog.Info("下一步:")
	kopylog.Info("  cd %s", renderedOutputDir)
	kopylog.Info("  # 根据项目类型进行初始化")
}
