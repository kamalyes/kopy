# kopy

> 通用项目模板生成工具 | Universal Project Template Generator

`kopy` 是一个类似 [Copier](https://github.com/copier-org/copier) 的通用模板生成工具，使用 Go 编写，跨平台支持（Windows / macOS / Linux）它不限于 Go 项目，可以生成任何语言的项目结构

## 特性

- 🌍 **跨平台** — Windows / macOS / Linux 原生支持
- 📦 **多源模板** — 支持本地目录和 Git 仓库（含分支指定 `@owner/repo@branch`）
- 🔄 **交互式收集** — 自动提示用户输入模板变量，支持 string/bool/int/select 四种类型
- 🎯 **命名转换** — 内置 snake/camel/pascal/kebab 转换（复用 [go-toolbox/stringx](https://github.com/kamalyes/go-toolbox)）
- ⚡ **自定义定界符** — 默认 `[[ ]]`，避免与 Go template 冲突
- 🛡️ **冲突检测** — 支持 skip / overwrite / ask 三种策略
- 🏃 **试运行模式** — `--dry-run` 预览生成结果
- 📋 **变量覆盖** — `--var Key=Value` 跳过交互
- 📂 **文件重命名** — 支持基于模板变量的文件/目录重命名规则
- 🔌 **可扩展** — 分层架构，每层可独立替换

## 安装

```bash
# 从源码构建
git clone https://github.com/kamalyes/kopy.git
cd kopy
make build

# 安装到 GOPATH/bin
make install

# 交叉编译
make cross
```

构建产物输出到 `bin/` 目录，交叉编译产物输出到 `dist/` 目录

## 快速开始

```bash
# 从本地模板生成项目
kopy new ./my-template

# 从 Git 仓库生成项目（@owner/repo 简写格式）
kopy new @kamalyes/gateway-template-service

# 指定分支（@owner/repo@branch 格式）
kopy new @kamalyes/gateway-template-service@v1.0

# 也支持完整 Git URL
kopy new https://github.com/kamalyes/gateway-template-service
kopy new https://github.com/kamalyes/gateway-template-service#v1.0

# 指定输出目录
kopy new ./my-template -o ./my-project

# 覆盖变量（跳过交互）
kopy new ./my-template --var ServiceName=Payment --var ModuleName=payment

# 强制覆盖已存在的文件
kopy new ./my-template --overwrite

# 跳过已存在的文件
kopy new ./my-template --skip

# 试运行（预览不写入）
kopy new ./my-template --dry-run

# 查看模板信息和问题列表
kopy list ./my-template
```

## 命令参考

### `kopy new <template-source>`

从模板生成新项目

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--output` | `-o` | `.` | 输出目录 |
| `--var` | `-v` | — | 模板变量覆盖，格式 `key=value`，可多次使用 |
| `--overwrite` | — | `false` | 覆盖已存在的文件 |
| `--skip` | — | `false` | 跳过已存在的文件 |
| `--dry-run` | — | `false` | 试运行模式，仅预览不写入 |

模板源格式：
- 本地目录：`./my-template` 或 `/path/to/template`
- Git 简写：`@owner/repo`（自动展开为 `https://github.com/owner/repo`）
- Git 带分支：`@owner/repo@branch`
- 完整 URL：`https://github.com/owner/repo` 或 `https://github.com/owner/repo#branch`

### `kopy list <template-source>`

列出模板信息和交互式问题列表

## 模板配置 (kopy.yaml)

每个模板项目根目录需要包含 `kopy.yaml` 配置文件，定义模板元信息、交互问题、忽略规则和重命名规则

```yaml
name: gateway-service
description: Go gRPC 微服务模板
version: 1.0.0
author: kamalyes
tags: [go, grpc, microservice]

# 模板定界符（默认 [[ ]]，避免与 Go template 冲突）
delimiters:
  - "[["
  - "]]"

# 交互式问题
questions:
  - name: ServiceName
    prompt: 服务名称（PascalCase）
    type: string
    default: MyService
    required: true

  - name: ModuleName
    prompt: 模块名称（小写）
    type: string
    default: "[[ .ServiceName | snake ]]"
    required: true

  - name: AuthorName
    prompt: 作者名称
    type: string
    default: kamalyes
    required: true

  - name: AuthorEmail
    prompt: 作者邮箱
    type: string
    default: example@qq.com
    required: true

  - name: GoVersion
    prompt: Go 版本
    type: select
    choices: ["1.22", "1.23", "1.24", "1.25"]
    default: "1.25"

  - name: UseRedis
    prompt: 是否使用 Redis
    type: bool
    default: "true"

  - name: Port
    prompt: 服务端口
    type: int
    default: "8080"

# 忽略的文件模式
ignore:
  - "kopy.yaml"
  - ".git/**"
  - "*.md"

# 文件重命名规则
rename:
  - pattern: "__service_name__"
    value: "[[ .ServiceName | snake ]]"

# 生命周期钩子
hooks:
  before:
    - echo "开始生成项目..."
  after:
    - go mod tidy
```

### 问题类型说明

| 类型 | 说明 | 输入方式 |
|------|------|----------|
| `string` | 文本输入 | 直接输入，回车确认 |
| `bool` | 布尔选择 | 输入 y/yes/true 或 n/no/false |
| `int` | 整数输入 | 输入数字，自动校验 |
| `select` | 下拉选择 | 输入选项编号 |

### 默认值模板语法

问题的 `default` 字段支持模板语法，可以引用之前已收集的变量：

```yaml
questions:
  - name: ServiceName
    type: string
    default: MyService
  - name: ModuleName
    type: string
    default: "[[ .ServiceName | snake ]]"  # 引用 ServiceName 并转为蛇形命名
```

### 内置变量

除了 `kopy.yaml` 中定义的问题变量外，引擎还会自动注入以下内置变量：

| 变量 | 说明 | 示例值 |
|------|------|--------|
| `Year` | 当前年份 | `2026` |

## 模板语法

使用 `[[ ]]` 作为定界符（可在 kopy.yaml 中自定义），所有模板文件的内容和文件名都会被渲染：

```go
// 文件名: [[ .ServiceName | snake ]]_service.go
package [[ .ModuleName | lower ]]

const ServiceName = "[[ .ServiceName ]]"
const ServicePort = [[ .Port ]]
```

文件头示例（动态注入作者信息）：

```go
/*
 * @Author: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @Date: [[ .Year ]]-04-19 00:00:00
 * @LastEditors: [[ .AuthorName ]] [[ .AuthorEmail ]]
 * @LastEditTime: [[ .Year ]]-04-19 00:00:00
 * @FilePath: \[[ .ServiceName | snake ]]\service.go
 * @Description: [[ .ServiceName ]] 服务实现
 *
 * Copyright (c) [[ .Year ]] by [[ .AuthorName ]], All Rights Reserved.
 */
```

### 内置函数

| 函数 | 说明 | 示例 |
|------|------|------|
| `snake` | 蛇形命名 | `MyService` → `my_service` |
| `camel` | 驼峰命名 | `my_service` → `myService` |
| `pascal` | 帕斯卡命名 | `my_service` → `MyService` |
| `kebab` | 短横线命名 | `MyService` → `my-service` |
| `lower` | 转小写 | `Hello` → `hello` |
| `upper` | 转大写 | `hello` → `HELLO` |
| `title` | 首字母大写 | `hello world` → `Hello World` |
| `trim` | 去除首尾空格 | ` hello ` → `hello` |
| `prefix` | 添加前缀 | `name` → `prefix_name` |
| `suffix` | 添加后缀 | `name` → `name_suffix` |
| `replace` | 替换字符串 | `hello world` → `hello_go` |
| `default` | 默认值 | `""` → `fallback` |

命名转换函数复用 [go-toolbox/stringx](https://github.com/kamalyes/go-toolbox) 实现

## 项目结构

```
kopy/
├── bootstrap/          # CLI 命令层
│   ├── root.go         # 根命令定义 + 全局参数
│   ├── new.go          # new 子命令 - 从模板生成新项目
│   ├── list.go         # list 子命令 - 列出模板信息
│   └── helpers.go      # 命令辅助函数 - 参数解析、冲突策略
├── config/             # 配置解析层
│   ├── types.go        # 配置类型定义 - 常量、枚举、结构体
│   ├── loader.go       # 配置加载 - 从目录读取并解析 kopy.yaml
│   ├── validate.go     # 配置校验 - 校验模板配置合法性
│   └── helper.go       # 配置辅助函数
├── engine/             # 渲染引擎层
│   ├── engine.go       # 引擎结构体定义与构造
│   ├── render.go       # 渲染方法 - 字符串、内容、路径渲染
│   └── funcs.go        # 内置模板函数 - 命名转换、字符串操作
├── prompt/             # 交互层
│   ├── prompt.go       # 提示器结构体与变量收集
│   ├── input.go        # 各类型输入方法 - string/bool/int/select
│   └── helper.go       # 辅助方法 - 默认值渲染
├── source/             # 模板源层
│   ├── source.go       # Source 接口与工厂 - 自动判断本地/Git源
│   ├── local.go        # 本地目录模板源实现
│   └── git.go          # Git 仓库模板源 - 克隆并查找模板目录
├── walker/             # 遍历层
│   ├── walker.go       # 模板目录递归遍历
│   ├── rename.go       # 文件重命名规则
│   └── ignore.go       # 忽略规则匹配
├── writer/             # 输出层
│   ├── writer.go       # 写入器结构体、Option 模式、冲突策略
│   ├── write.go        # 文件写入方法 - 渲染并写入，支持冲突检测
│   └── helper.go       # 辅助函数 - 冲突询问、可执行文件检测
├── main.go             # 入口
├── Makefile            # 构建脚本
└── go.mod
```

## 架构设计

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ bootstrap/     │────▶│  source/ │────▶│ config/  │────▶│ prompt/  │
│  CLI层   │     │  源解析   │     │  配置解析  │     │  变量收集  │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                                          │
                                                          ▼
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│ writer/  │◀────│ walker/  │◀────│ engine/  │◀────│  vars    │
│  文件输出  │     │  目录遍历  │     │  模板渲染  │     │  变量集   │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
```

### 执行流程

1. **源解析** (`source/`) — 根据输入判断本地路径或 Git 仓库，Git 仓库自动克隆到临时目录
2. **配置加载** (`config/`) — 从模板目录读取 `kopy.yaml`，解析并校验配置
3. **变量收集** (`prompt/`) — 交互式收集用户输入，`--var` 参数可跳过交互
4. **模板渲染** (`engine/`) — 使用 Go `text/template` 引擎渲染文件内容和路径
5. **目录遍历** (`walker/`) — 递归遍历模板目录，应用忽略规则和重命名规则
6. **文件写入** (`writer/`) — 渲染后的内容写入目标目录，支持冲突检测和试运行

### 扩展点

每一层职责单一，可独立替换：
- **source/** — 扩展新的模板源（如 S3、OCI Registry）
- **engine/** — 扩展新的模板函数或定界符
- **prompt/** — 替换为 TUI 界面
- **writer/** — 扩展新的冲突策略

## 忽略规则

`kopy.yaml` 中的 `ignore` 字段支持多种匹配模式：

| 模式 | 说明 | 示例 |
|------|------|------|
| 精确匹配 | 完全匹配路径 | `kopy.yaml` |
| 前缀匹配 | 匹配目录前缀 | `.git/` |
| 后缀匹配 | 匹配文件扩展名 | `*.md` |
| 通配符匹配 | glob 模式 | `docs/**/*.md` |
| 路径段匹配 | 匹配路径中任意段 | `__pycache__` |

默认忽略规则：`kopy.yaml` 和 `.git/**`

## 文件重命名

`rename` 字段定义文件/目录名的替换规则，支持模板语法：

```yaml
rename:
  - pattern: "__service_name__"
    value: "[[ .ServiceName | snake ]]"
  - pattern: "__module_name__"
    value: "[[ .ModuleName | lower ]]"
```

模板目录中的 `__service_name__` 目录或文件名会被替换为渲染后的值

## 冲突处理策略

当目标路径已存在文件时，支持三种处理策略：

| 策略 | 参数 | 说明 |
|------|------|------|
| `ask` | 默认 | 交互式询问是否覆盖 |
| `overwrite` | `--overwrite` | 强制覆盖所有已存在文件 |
| `skip` | `--skip` | 跳过所有已存在文件 |

## 可执行文件检测

写入器会自动检测文件内容是否以 shebang（`#!`）开头，如果是则设置可执行权限（0755），否则使用普通文件权限（0644）

## 官方模板

| 模板 | 说明 |
|------|------|
| [gateway-template-proto](https://github.com/kamalyes/gateway-template-proto) | Protobuf 协议中心模板 |
| [gateway-template-service](https://github.com/kamalyes/gateway-template-service) | Go gRPC 微服务模板 |

## 依赖

| 依赖 | 说明 |
|------|------|
| [go-toolbox](https://github.com/kamalyes/go-toolbox) | 命名转换（stringx）、数学工具（mathx）等 |
| [cobra](https://github.com/spf13/cobra) | CLI 框架 |
| [yaml.v3](https://gopkg.in/yaml.v3) | YAML 解析 |

## License

Copyright (c) 2026 by kamalyes, All Rights Reserved.
