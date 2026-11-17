# 写模板 | Authoring a template

> 怎么写一个 kopy 模板：`kopy.yaml` 的每个字段、问题类型、默认值派生、模板语法与内置函数。
> 规则类内容（忽略与重命名）见 [RULES.md](./RULES.md)，生成后初始化见 [HOOKS.md](./HOOKS.md)。

> [EN] How to author a kopy template: every `kopy.yaml` field, question types, derived
> defaults, template syntax and built-in functions. Ignore and rename rules live in
> [RULES.md](./RULES.md); post-generate initialization in [HOOKS.md](./HOOKS.md).

## ⚙️ 模板配置 | Template configuration (`kopy.yaml`)

```yaml
name: my-template
description: 模板描述 | Template description
version: 1.0.0
author: kamalyes
tags: [go, grpc]

# 模板定界符（默认 [[ ]]）| Delimiters, default [[ ]]
# 模板里带 shell 脚本时务必换掉，见「重要行为」第 1 条
delimiters:
  - "<%"
  - "%>"

# 交互式问题 | Questions
questions:
  - name: ServiceName
    prompt: "服务名称 | Service name"
    type: string
    default: MyService
    required: true

  - name: ModuleName
    prompt: "模块名 | Module name"
    type: string
    # 引用前面已答的变量，顺序敏感
    default: "<% .ServiceName | snake %>"

  - name: GoVersion
    prompt: "Go 版本 | Go version"
    type: select
    choices: ["1.23", "1.24", "1.25"]
    default: "1.25"

  - name: UseRedis
    prompt: "是否使用 Redis | Use Redis"
    type: bool
    default: "true"

  - name: Port
    prompt: "服务端口 | Port"
    type: int
    default: "8080"

# 忽略的文件模式（匹配渲染后的路径，模式本身也会渲染）
# [EN] Ignore patterns; matched against the rendered path and themselves rendered
ignore:
  - "kopy.yaml"
  - ".git/"
  - ".git/**"
  - '<% if ne .UseRedis "true" %>cache/redis.go<% end %>'

# 文件重命名规则 | Rename rules
rename:
  - pattern: "__service_name__"
    value: "<% .ServiceName | snake %>"

# 生成后执行的初始化命令 | Post-generate commands
hooks:
  after:
    - bash scripts/init.sh
```

> `kopy.yaml` 由内置的 YAML 子集解析器读取：支持块级映射/序列、flow 序列/映射、引号、注释、`|` / `>` 块标量；不支持锚点、别名、多文档和多行 flow 集合。写模板时按上面的示例使用即可，无需引入额外 YAML 库。
>
> [EN] `kopy.yaml` is read by a built-in YAML subset parser. Block mappings/sequences,
> flow sequences/mappings, quotes, comments and `|` / `>` block scalars are supported;
> anchors, aliases, multiple documents and multi-line flow collections are not.

### 问题类型 | Question types

| 类型 \| Type | 说明 \| Description | 输入方式 \| Input |
| --- | --- | --- |
| `string` | 文本 | 直接输入，回车确认 |
| `bool` | 布尔 | `y`/`yes`/`true` 或 `n`/`no`/`false` |
| `int` | 整数 | 输入数字，自动校验 |
| `select` | 单选 | 输入选项编号 |

所有变量值在模板里**一律是字符串**，判断布尔要写 `eq .UseRedis "true"`。

[EN] Every value is a string inside templates, so booleans are compared via
`eq .UseRedis "true"`.

### 内置变量 | Built-in variables

| 变量 \| Variable | 说明 \| Description | 示例 \| Example |
| --- | --- | --- |
| `Year` | 当前年份，问题的 `default` 里也可用 \| current year, usable in defaults too | `2026` |

## 🧬 默认值可以引用前面已答的变量 | Defaults may reference earlier answers

每答完一题就写回引擎，所以后面问题的 `default` 能引用前面的答案，过滤器也可用：

```yaml
questions:
  - name: OrgPrefix
    type: string
    default: demo
  - name: ProjectName
    type: string
    default: "[[ .OrgPrefix ]]-share-proto"     # → demo-share-proto
  - name: DBName
    type: string
    default: "[[ .OrgPrefix | snake ]]_service" # → demo_service
```

顺序敏感：被引用的问题必须排在前面。前向引用会渲染出 `<no value>`，kopy 会给出告警而不是静默放过。

[EN] Each answer is written back to the engine as it is collected, so later defaults can
reference earlier answers, filters included. Order matters: a referenced question must come
first. A forward reference renders as `<no value>` and kopy warns rather than passing silently.

## 🔤 模板语法 | Template syntax

文件的**内容**和**路径**都会过模板引擎：

```go
// 文件名 | file name: <% .ServiceName | snake %>_service.go
package <% .ModuleName | lower %>

const ServiceName = "<% .ServiceName %>"
const ServicePort = <% .Port %>
```

文件头示例 | file header example：

```go
/*
 * @Author: <% .AuthorName %> <% .AuthorEmail %>
 * @Date: <% nowFmt "2006-01-02 15:04:05" %>
 * @FilePath: \<% .ProjectName %>\service\<% .ServiceName | snake %>.go
 * @Description: <% .ServiceName %> 服务实现
 *
 * Copyright (c) <% .Year %> by <% .AuthorName %>, All Rights Reserved.
 */
```

### 内置函数 | Built-in functions

| 函数 \| Function | 说明 \| Description | 示例 \| Example |
| --- | --- | --- |
| `snake` | 蛇形 | `MyService` → `my_service` |
| `camel` | 驼峰 | `my_service` → `myService` |
| `pascal` | 帕斯卡 | `my_service` → `MyService` |
| `kebab` | 短横线 | `MyService` → `my-service` |
| `lower` | 转小写 | `Hello` → `hello` |
| `upper` | 转大写 | `hello` → `HELLO` |
| `title` | 首字母大写 | `hello world` → `Hello World` |
| `trim` | 去首尾空格 | ` hello ` → `hello` |
| `prefix` | 加前缀 | `name` → `prefix_name` |
| `suffix` | 加后缀 | `name` → `name_suffix` |
| `replace` | 替换 | `hello world` → `hello_go` |
| `now` | 当前时间（`time.Time`）| `<% now %>` |
| `date` | 按格式渲染时间 | `<% now \| date "2006-01-02 15:04:05" %>` |
| `nowFmt` | 按格式渲染当前时间 | `<% nowFmt "2006-01-02" %>` |
| `default` | 默认值 | `""` → `fallback` |

命名转换为内置实现（纯标准库），无外部依赖。
时间格式用 Go Layout：`2006-01-02 15:04:05` 对应 `yyyy-mm-dd HH:MM:SS`。

## ⚠️ 定界符会和 shell 脚本打架 | Delimiters collide with shell scripts

默认定界符是 `[[ ]]`，而 bash 的条件测试语法本身就是双方括号（`if [[ -t 1 ]]`），字符类也是（`[[:space:]]`）。模板里只要带 shell 脚本，用默认定界符会直接渲染失败：

```
template: inline:21: bad number syntax: "-t"
```

模板里需要同步生产仓库的 shell 脚本时，把定界符改掉：

```yaml
delimiters:
  - "<%"
  - "%>"
```

这样 `.sh` 可以逐字节搬进来，不必为了模板引擎把 `[[ ]]` 改写成 `[ ]`（改写会悄悄改变语义，比如 `=~` 和 `&&`）。

[EN] The default `[[ ]]` collides with bash's own test syntax and character classes, so a
template containing shell scripts fails to render. Switch the delimiters to `<% %>` and shell
scripts sync in byte for byte, instead of rewriting `[[ ]]` into `[ ]` — a rewrite that
silently changes semantics.
