# 生命周期钩子 | Lifecycle hooks

> 生成完成后自动跑初始化命令，让 `kopy new` 一步到可运行的项目。
> 这一步会在使用者机器上执行模板作者写的命令，所以默认先询问——取向与实现见下。

> [EN] Run the template's initialization commands after generating, so `kopy new` lands a
> runnable project in one step. Because this executes the template author's commands on the
> user's machine, kopy asks first.

## 🪝 声明与执行 | Declaring and running

```yaml
hooks:
  after:
    - bash scripts/init.sh
```

命令在**生成出来的项目目录内**按顺序执行，输出实时透传（初始化动辄几十秒，攒到最后再打印
等于黑屏等待）。命令本身会先渲染，所以可以引用变量。执行策略与安全取向见「重要行为」第 4 条。

[EN] Commands run inside the generated project directory in order, with output streamed
through. They are rendered first, so they may reference variables.

## 🛡️ 为什么默认先问 | Why it asks first

模板源可以是任意 Git 地址，而 hooks 就是任意命令执行——静默执行等于 `kopy new <url>` 即可在使用者机器上跑代码。所以默认先把**渲染后**的命令逐条打印出来，再问一次，回车即否：

```
🪝 模板声明了 1 条初始化命令（hooks.after），将在 /path/to/out 内执行：
  1) bash scripts/init.sh
  是否执行？这会在你的机器上运行上述命令 [y/N]:
```

| 取值 \| Value | 行为 \| Behaviour |
| --- | --- |
| `--hooks ask` | 列出命令后询问，回车即否（默认）|
| `--hooks always` | 不询问直接执行，CI 用 |
| `--hooks never` | 完全跳过 |
| `--dry-run` | 一律不执行 |

任一条失败即停止，并明确说明项目文件已生成，不让人误以为整次生成失败。
`hooks.before` 当前不执行（渲染前还没有项目目录，没有合适的工作目录），声明了会给告警。

[EN] A template source can be any Git URL and hooks are arbitrary command execution, so
running them silently would mean `kopy new <url>` executes code on the user's machine. kopy
therefore prints the rendered commands and asks, defaulting to no. Use `--hooks always` in
CI. Execution stops at the first failure and states that the files were already generated.
`hooks.before` is not executed yet and only warns.
