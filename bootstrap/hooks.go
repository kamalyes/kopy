/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-20 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-20 00:00:00
 * @FilePath: \kopy\bootstrap\hooks.go
 * @Description: 生命周期钩子执行 - 渲染并在生成项目内执行 hooks.after
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/kamalyes/kopy/config"
	"github.com/kamalyes/kopy/engine"
	kopylog "github.com/kamalyes/kopy/log"
	"github.com/kamalyes/kopy/prompt"
)

// 钩子执行模式
const (
	// HookModeAsk 列出命令后询问一次，默认否
	HookModeAsk = "ask"
	// HookModeAlways 不询问直接执行，供 CI 使用
	HookModeAlways = "always"
	// HookModeNever 完全跳过
	HookModeNever = "never"
)

// runAfterHooks 渲染并在生成项目目录内依次执行 hooks.after
//
// 为什么默认是「问一次」而不是静默执行：模板源可以是任意 Git 地址，
// hooks 就是任意命令执行。静默执行等于 `kopy new <url>` 即可在使用者机器上跑代码，
// 所以默认先把渲染后的命令逐条打印出来，再问一次（回车即否）。
// CI 用 --hooks always，完全不要用 --hooks never；--dry-run 一律不执行。
//
// 命令按 shell 字符串执行（unix 用 sh -c，Windows 用 cmd /C），因为模板里写的是
// `bash scripts/init.sh`、`go mod tidy` 这类带参数的命令行，不是 argv 数组。
// 工作目录是生成出来的项目根目录，这样 hooks 才能直接引用项目内的相对路径。
// 任一条失败即停止并返回错误：后续步骤通常依赖前一步的产物，继续跑只会刷更多错误。
func runAfterHooks(cfg *config.TemplateConfig, eng *engine.Engine, p *prompt.Prompter,
	projectDir, mode string, dryRun bool) error {
	if len(cfg.Hooks.Before) > 0 {
		// before 钩子目前不执行：渲染前还没有项目目录，没有合适的工作目录与语义。
		// 明确告知而不是静默忽略，免得模板作者以为它跑过了
		kopylog.Warn("⚠️  hooks.before 声明了 %d 条命令，当前版本不执行（渲染前尚无项目目录）",
			len(cfg.Hooks.Before))
	}

	if len(cfg.Hooks.After) == 0 {
		return nil
	}

	if dryRun {
		kopylog.Info("🪝 试运行模式：hooks.after 的 %d 条命令不会执行", len(cfg.Hooks.After))
		return nil
	}

	if mode == HookModeNever {
		kopylog.Info("🪝 按 --hooks never 跳过 hooks.after 的 %d 条命令", len(cfg.Hooks.After))
		return nil
	}

	// 先整体渲染，拿到真正要执行的命令再决定是否执行：
	// 让使用者看到的就是将要跑的那一行，而不是带占位符的原文
	commands := make([]string, 0, len(cfg.Hooks.After))
	for _, raw := range cfg.Hooks.After {
		rendered, err := eng.RenderString(raw)
		if err != nil {
			return fmt.Errorf("渲染钩子命令 %q 失败: %w", raw, err)
		}
		rendered = strings.TrimSpace(rendered)
		if rendered == "" {
			continue
		}
		commands = append(commands, rendered)
	}
	if len(commands) == 0 {
		return nil
	}

	kopylog.Info("🪝 模板声明了 %d 条初始化命令（hooks.after），将在 %s 内执行：",
		len(commands), projectDir)
	for i, cmd := range commands {
		kopylog.Info("  %d) %s", i+1, cmd)
	}

	if mode == HookModeAsk {
		// 默认否是刻意的：这一步会在使用者机器上运行模板作者写的命令，
		// 不该因为随手回车就发生
		if !p.Confirm("  是否执行？这会在你的机器上运行上述命令") {
			kopylog.Info("🪝 已跳过。需要时可手动执行上面的命令，或重跑时加 --hooks always")
			return nil
		}
	}

	for i, cmd := range commands {
		kopylog.Info("🪝 [%d/%d] %s", i+1, len(commands), cmd)
		if err := runShell(cmd, projectDir); err != nil {
			return fmt.Errorf("第 %d 条钩子命令失败（%s）: %w", i+1, cmd, err)
		}
	}
	kopylog.Info("✅ 初始化命令执行完成")
	return nil
}

// runShell 在指定目录内以 shell 执行一条命令，输出直接透传
func runShell(command, dir string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.Command(shell, "-c", command)
	}
	cmd.Dir = dir
	// 透传而不是收集：初始化动辄几十秒（装 protoc、下依赖），
	// 使用者需要实时看到进度，攒到最后再打印等于黑屏等待
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
