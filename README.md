# kopy - 通用项目模板生成工具

> `kopy` 是一个用 Go 写的通用项目模板生成工具，跨平台（Windows / macOS / Linux），不限语言：模板里放什么语言的工程，它就能生成什么。交互式收集变量、渲染文件内容与路径、按变量条件性排除文件、生成后自动执行初始化命令。

> [EN] `kopy` is a universal project template generator written in Go. It runs natively on Windows, macOS and Linux and is language-agnostic: whatever stack the template holds is what it generates. It collects variables interactively, renders both file contents and paths, excludes files per variable, and runs the template's initialization commands after generating.

[![stable](https://img.shields.io/badge/stable-stable-green.svg)](https://github.com/kamalyes/kopy)
[![license](https://img.shields.io/github/license/kamalyes/kopy)]()
[![release](https://img.shields.io/github/v/release/kamalyes/kopy)]()
[![commit](https://img.shields.io/github/last-commit/kamalyes/kopy)]()
[![issues](https://img.shields.io/github/issues/kamalyes/kopy)]()
[![pull](https://img.shields.io/github/issues-pr/kamalyes/kopy)]()
[![fork](https://img.shields.io/github/forks/kamalyes/kopy)]()
[![star](https://img.shields.io/github/stars/kamalyes/kopy)]()
[![go](https://img.shields.io/github/go-mod/go-version/kamalyes/kopy)]()
[![size](https://img.shields.io/github/repo-size/kamalyes/kopy)]()
[![Go Reference](https://pkg.go.dev/badge/github.com/kamalyes/kopy?status.svg)](https://pkg.go.dev/github.com/kamalyes/kopy?tab=doc)

## 📚 文档导航 | Documentation

### 📖 本仓库文档 | In this repository

| 文档 \| Document | 内容 \| Contents |
| --- | --- |
| [写模板 \| Authoring](./docs/AUTHORING.md) | `kopy.yaml` 每个字段、问题类型、默认值派生、模板语法与内置函数、定界符与 shell 脚本的冲突 |
| [忽略与重命名 \| Ignore and rename](./docs/RULES.md) | `ignore` / `rename` 规则、跨层级匹配陷阱、按变量条件性排除文件 |
| [生命周期钩子 \| Hooks](./docs/HOOKS.md) | `hooks.after` 的执行与安全取向、`--hooks` 取值、失败处理 |
| [架构与执行流程 \| Architecture](./docs/ARCHITECTURE.md) | 六层架构、七步流程、扩展点、冲突策略 |

### 🔗 外部链接 | External

- [🏠 项目主页 \| Home](https://github.com/kamalyes/kopy)
- [📖 API 文档 \| API reference](https://pkg.go.dev/github.com/kamalyes/kopy)
- [🐛 问题反馈 \| Issues](https://github.com/kamalyes/kopy/issues)

## 🚀 核心特性 | Features

- 🌍 **跨平台 | Cross-platform** — Windows / macOS / Linux 原生支持
- 📦 **多源模板 | Multiple sources** — 本地目录与 Git 仓库，支持 `@owner/repo@branch` 简写
- 🔄 **交互式收集 | Interactive prompts** — string / bool / int / select 四种类型，`--yes` 可全程零交互
- 🧬 **默认值可派生 | Derived defaults** — 后面问题的 `default` 能引用前面已答的变量
- 🎯 **命名转换 | Name conversion** — 内置 snake / camel / pascal / kebab（复用 [go-toolbox/stringx](https://github.com/kamalyes/go-toolbox)）
- ⚡ **自定义定界符 | Custom delimiters** — 默认 `[[ ]]`，可改成 `<% %>` 等
- 🚫 **条件性排除 | Conditional exclusion** — `ignore` 模式本身会被渲染，可按变量决定某个文件或目录是否生成
- 📂 **文件重命名 | Renaming** — 基于模板变量重写文件与目录名
- 🪝 **生成后初始化 | Post-generate hooks** — `hooks.after` 会被执行，默认先列出命令再询问
- 🛡️ **冲突策略 | Conflict strategy** — skip / overwrite / ask
- 🏃 **试运行 | Dry run** — `--dry-run` 只预览不写入
- 🔌 **分层可替换 | Layered** — source / config / engine / prompt / walker / writer 各层职责单一

## ⚠️ 四个必读行为 | Four behaviours to know first

写模板前先过一遍这四条，它们决定了模板能怎么写；每条都有踩过的实际后果。

[EN] These four decide what a template can express, and each has a concrete failure behind it.

| # | 行为 \| Behaviour | 不知道的后果 \| What goes wrong | 详情 \| Details |
| --- | --- | --- | --- |
| 1 | 默认定界符 `[[ ]]` 与 bash 的 `if [[ ]]`、`[[:space:]]` 冲突 | 模板一带 shell 脚本就渲染失败：`bad number syntax: "-t"` | [AUTHORING](./docs/AUTHORING.md) |
| 2 | `ignore` 判定的是**渲染后**的路径，模式本身也会渲染 | 不知道就只能全量生成后手删；写了 `*.md` 会把生成项目的全部文档删光 | [RULES](./docs/RULES.md) |
| 3 | 问题的 `default` 可以引用**前面已答**的变量，顺序敏感 | 前向引用渲染出 `<no value>`，会被当成答案写进每个文件 | [AUTHORING](./docs/AUTHORING.md) |
| 4 | `hooks.after` 会执行，但默认先列出命令再询问（回车即否）| CI 里不加 `--hooks always` 会卡在询问 | [HOOKS](./docs/HOOKS.md) |

## 📦 安装 | Installation

```bash
# 指定版本安装到 GOPATH/bin（推荐）| pinned install
go install github.com/kamalyes/kopy@latest

# 从源码构建 | build from source
git clone https://github.com/kamalyes/kopy.git && cd kopy
go build -o kopy .

# 安装到 GOPATH/bin | install into GOPATH/bin
go install .

# 交叉编译 | cross compile
GOOS=windows GOARCH=amd64 go build -o kopy.exe .
GOOS=linux   GOARCH=arm64 go build -o kopy-linux-arm64 .
```

## 🚀 快速开始 | Quick start

```bash
# 从本地模板生成 | from a local template
kopy new ./my-template -o ./my-project

# 从 Git 仓库生成（@owner/repo 简写）| from Git, shorthand
kopy new @kamalyes/kronos-service-template -o ./my-service

# 指定分支 | pick a branch
kopy new @kamalyes/kronos-service-template@v1.0 -o ./my-service

# 完整 URL 同样支持 | full URL works too
kopy new https://github.com/kamalyes/kronos-service-template -o ./my-service

# 零交互，全部用默认值 | zero prompts, all defaults
kopy new ./my-template -o ./my-project -y

# 零交互 + 指定几个变量（--var 优先，派生默认值会跟着级联）
# [EN] zero prompts plus a few overrides; --var wins and derived defaults cascade
kopy new ./my-template -o ./my-project -y -v ServiceType=payment

# 生成后自动执行初始化命令，不询问（CI）| run hooks without asking
kopy new ./my-template -o ./my-project -y --hooks always

# 试运行，只预览不写入 | preview only
kopy new ./my-template -o ./my-project --dry-run

# 查看模板信息与问题列表 | inspect the template
kopy list ./my-template
```

## 📖 命令参考 | Command reference

### `kopy new <template-source>`

从模板生成新项目 | Generate a project from a template

| 参数 \| Flag | 简写 | 默认值 | 说明 \| Description |
| --- | --- | --- | --- |
| `--output` | `-o` | `.` | 输出目录 \| Output directory |
| `--var` | `-v` | — | 变量覆盖，`key=value`，可多次使用 \| Variable override, repeatable |
| `--yes` | `-y` | `false` | 不询问，全部使用默认值（`--var` 仍优先）\| Skip all prompts and use defaults |
| `--hooks` | — | `ask` | `hooks.after` 执行方式：`ask` / `always` / `never` |
| `--overwrite` | — | `false` | 覆盖已存在的文件 \| Overwrite existing files |
| `--skip` | — | `false` | 跳过已存在的文件 \| Skip existing files |
| `--dry-run` | — | `false` | 试运行，仅预览不写入 \| Preview without writing |

`--yes` 下必填项没有默认值时会直接报错并提示用 `--var` 指定，不会静默生成一个缺值的项目；
`select` 的默认值不在候选里同样报错（模板定义有误）。

[EN] Under `--yes`, a required question without a default errors out and asks for `--var`
rather than silently producing a project with a missing value; a `select` default outside its
choices errors too.

模板源格式 | Template source formats：

- 本地目录 | local path：`./my-template`、`/path/to/template`
- Git 简写 | shorthand：`@owner/repo`（展开为 `https://github.com/owner/repo`）
- Git 带分支 | with branch：`@owner/repo@branch`
- 完整 URL | full URL：`https://github.com/owner/repo`、`https://github.com/owner/repo#branch`

### `kopy list <template-source>`

列出模板信息与交互式问题列表 | List the template's metadata and questions

## 📋 官方模板 | Official templates

| 模板 \| Template | 说明 \| Description |
| --- | --- |
| [kronos-share-proto-template](https://github.com/kamalyes/kronos-share-proto-template) | gRPC/Protobuf 协议中心：多模块 proto、代码生成链路、网关文档聚合与 CI |
| [kronos-service-template](https://github.com/kamalyes/kronos-service-template) | Go gRPC 微服务：接协议中心、pbmo 转换、sqlbuilder 仓储、go-bizerror 错误 |
| [openapi-ts-sdk-template](https://github.com/kamalyes/openapi-ts-sdk-template) | 从远程 OpenAPI 文档生成、可发布到 npm 的 TypeScript SDK |

## 📦 依赖 | Dependencies

| 依赖 \| Dependency | 说明 \| Description |
| --- | --- |
| [go-toolbox](https://github.com/kamalyes/go-toolbox) | 命名转换（stringx）、数学工具（mathx） |
| [go-logger](https://github.com/kamalyes/go-logger) | 日志输出 |
| [cobra](https://github.com/spf13/cobra) | CLI 框架 |
| [yaml.v3](https://gopkg.in/yaml.v3) | YAML 解析 |

## License

MIT
