# 忽略与重命名 | Ignore and rename

> `ignore` 与 `rename` 决定「哪些文件进入生成项目、叫什么名字」。
> 两者都在渲染之后生效，所以可以按变量条件性地排除文件——这是功能开关不留孤立文件的做法。
> 模板写法见 [AUTHORING.md](./AUTHORING.md)。

> [EN] `ignore` and `rename` decide which files reach the generated project and under what
> names. Both act after rendering, so files can be excluded per variable.

## 🎯 判定对象是渲染后的路径与模式 | Matched against the RENDERED path and patterns

`ignore` 的匹配发生在 rename 与渲染**之后**，并且模式本身也会先渲染一遍。所以「按变量决定某个目录是否生成」可以直接表达：

```yaml
questions:
  - name: GenerateExample
    type: bool
    default: "true"

ignore:
  # GenerateExample 为否时整组示例文件不生成
  - '[[ if ne .GenerateExample "true" ]]models/[[ .ExampleName | snake ]]_model.go[[ end ]]'
```

也可以配合 rename 落到一个跳过标记上：

```yaml
ignore:
  - "__kopy_skip__/"
  - "__kopy_skip__/**"
rename:
  - pattern: "__maybe__"
    value: '[[ if eq .Flavor "keep" ]]kept[[ else ]]__kopy_skip__[[ end ]]'
```

[EN] ignore runs after rename and rendering, and the patterns are rendered too, so a file or
directory can be excluded per variable — either directly in the pattern, or by renaming it
onto a skip marker.

> ⚠️ 不要写 `*.md`：它按后缀匹配**所有层级**的 markdown，会把生成项目的全部文档一起删掉。
> 要排除模板自身的说明文档，用 rename 落到一个独有标记再忽略那个标记。
>
> [EN] Never use `*.md`: it suffix-matches markdown at every level and wipes the generated
> project's entire documentation. To keep the template's own docs out, rename them onto a
> unique marker and ignore the marker.

## 🚫 忽略规则 | Ignore rules

`ignore` 支持多种匹配模式，判定对象是**渲染后**的路径：

| 模式 \| Pattern | 说明 \| Description | 示例 \| Example |
| --- | --- | --- |
| 精确匹配 | 完全匹配整条相对路径 | `kopy.yaml` |
| 前缀匹配 | 以 `/` 结尾，匹配目录前缀 | `.git/` |
| 后缀匹配 | 以 `*.` 开头，匹配扩展名（**所有层级**）| `*.log` |
| 通配符 | glob 模式 | `docs/**` |
| 路径段匹配 | 匹配路径中任意一段 | `__pycache__` |

注意后缀匹配与路径段匹配都是**跨层级**的：`*.md` 会命中 `README.md` 也会命中
`locales/README.md`，`README.md` 作为模式同样会因「路径段匹配」命中子目录里的同名文件。
要只排除根目录那一份，用 rename 落到一个独有标记：

```yaml
ignore:
  - "__kopy_template_doc__"
rename:
  - pattern: "README.md"
    value: "__kopy_template_doc__"
```

[EN] Both suffix and path-segment matching cross directory levels, so `*.md` hits
`locales/README.md` too, and `README.md` as a pattern hits same-named files in
subdirectories. To exclude only the root one, rename it onto a unique marker.

## 📂 文件重命名 | Renaming

`rename` 按顺序对同一条路径字符串做 `ReplaceAll`，值支持模板语法：

```yaml
rename:
  - pattern: ".yml.tpl"
    value: ".yml"
  - pattern: "__service_name__"
    value: "<% .ServiceName | snake %>"
```

顺序敏感，且要留意**子串包含关系**：`__module2__` 含有 `__module__` 吗？不含，
所以两条规则互不干扰；反之若 `__mod__` 与 `__mod_x__` 并存，前者会先把后者改坏。

工作流文件常带 `.yml.tpl` 后缀再由 rename 去掉：GitHub Actions 会解析模板仓库自己的
`.github/workflows/*.yml`，而里面的占位符尚未渲染（`branches: - <% .BaseBranch %>` 会被
YAML 读成嵌套序列），模板仓库会一直挂着无效工作流的红叉。

[EN] rename rules are sequential `ReplaceAll` calls, so watch for substring containment.
Workflow files usually carry a `.yml.tpl` suffix that rename strips, because GitHub Actions
parses the template repository's own workflows where placeholders are still unrendered.
