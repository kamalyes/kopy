/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\bootstrap\cli.go
 * @Description: 轻量命令与参数解析 - 基于标准库 flag，替代 cobra
 *
 * 支持子命令、全局 flag、flag 与位置参数混排、短选项聚合（-yh）、
 * --flag=value、-ovalue 等常见形式，帮助输出保持与 cobra 相近的风格。
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bootstrap

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// ==================== 命令 ====================

// Command 命令定义
type Command struct {
	Use   string
	Short string
	Long  string
	Args  func(cmd *Command, args []string) error
	Run   func(cmd *Command, args []string)

	// Out/Err 仅用于测试或嵌入时重定向输出，为空时使用标准输出/错误
	Out io.Writer
	Err io.Writer

	flags    *FlagSet
	parent   *Command
	children []*Command
}

// AddCommand 添加子命令
func (c *Command) AddCommand(children ...*Command) {
	for _, child := range children {
		child.parent = c
		c.children = append(c.children, child)
	}
}

func (c *Command) root() *Command {
	if c.parent == nil {
		return c
	}
	return c.parent.root()
}

func (c *Command) commandName() string {
	if i := strings.IndexByte(c.Use, ' '); i >= 0 {
		return c.Use[:i]
	}
	return c.Use
}

func (c *Command) fullUse() string {
	if c.parent == nil {
		return c.Use
	}
	return c.parent.fullUse() + " " + c.Use
}

func (c *Command) stdout() io.Writer {
	if c.Out != nil {
		return c.Out
	}
	if c.parent != nil {
		return c.parent.stdout()
	}
	return os.Stdout
}

func (c *Command) stderr() io.Writer {
	if c.Err != nil {
		return c.Err
	}
	if c.parent != nil {
		return c.parent.stderr()
	}
	return os.Stderr
}

// PersistentFlags 返回命令的全局 flag 集合
func (c *Command) PersistentFlags() *FlagSet {
	if c.flags == nil {
		c.flags = newFlagSet(c.commandName())
	}
	return c.flags
}

// Flags 返回当前命令可见的 flag 集合
func (c *Command) Flags() *FlagSet {
	return c.root().PersistentFlags()
}

func (c *Command) findChild(name string) *Command {
	for _, child := range c.children {
		if child.commandName() == name {
			return child
		}
	}
	return nil
}

// Execute 使用 os.Args 执行命令
func (c *Command) Execute() error {
	return c.ExecuteArgs(os.Args[1:])
}

// ExecuteArgs 使用给定参数执行命令
func (c *Command) ExecuteArgs(args []string) error {
	root := c.root()
	flags := root.PersistentFlags()

	positional, err := flags.parse(args)
	if err != nil {
		return c.fail(err, c)
	}
	return c.dispatch(positional, flags)
}

func (c *Command) dispatch(args []string, flags *FlagSet) error {
	help, _ := flags.GetBool("help")

	if len(args) == 0 {
		if c.Run == nil || help {
			c.printHelp(c.stdout())
			return nil
		}
		if c.Args != nil {
			if err := c.Args(c, nil); err != nil {
				return c.fail(err, c)
			}
		}
		c.Run(c, nil)
		return nil
	}

	if args[0] == "help" {
		if len(args) == 1 {
			c.printHelp(c.stdout())
			return nil
		}
		target := c.findChild(args[1])
		if target == nil {
			return c.fail(fmt.Errorf("unknown help topic %q", args[1]), nil)
		}
		target.printHelp(target.stdout())
		return nil
	}

	child := c.findChild(args[0])
	if child == nil {
		return c.fail(fmt.Errorf("unknown command %q for %q", args[0], c.commandName()), nil)
	}

	if help {
		child.printHelp(child.stdout())
		return nil
	}

	childArgs := args[1:]
	if child.Args != nil {
		if err := child.Args(child, childArgs); err != nil {
			return c.fail(err, child)
		}
	}
	if child.Run == nil {
		child.printHelp(child.stdout())
		return nil
	}

	child.Run(child, childArgs)
	return nil
}

func (c *Command) fail(err error, cmd *Command) error {
	writef(c.stderr(), "Error: %v\n", err)
	if cmd != nil {
		writeln(c.stderr())
		cmd.printUsage(c.stderr())
		cmd.printFlags(c.stderr())
	} else {
		writef(c.stderr(), "Run '%s --help' for usage.\n", c.commandName())
	}
	return err
}

// ExactArgs 校验位置参数个数
func ExactArgs(n int) func(cmd *Command, args []string) error {
	return func(cmd *Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("accepts %d arg(s), received %d", n, len(args))
		}
		return nil
	}
}

// ==================== 帮助输出 ====================

func writef(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func writeln(w io.Writer, args ...any) {
	_, _ = fmt.Fprintln(w, args...)
}

func (c *Command) printHelp(w io.Writer) {
	description := c.Long
	if description == "" {
		description = c.Short
	}
	if description != "" {
		writeln(w, description)
		writeln(w)
	}

	c.printUsage(w)
	c.printFlags(w)

	if len(c.children) > 0 {
		writeln(w)
		writef(w, "Use %q for more information about a command.\n",
			c.commandName()+" [command] --help")
	}
}

func (c *Command) printUsage(w io.Writer) {
	writeln(w, "Usage:")
	if len(c.children) > 0 {
		writef(w, "  %s [command]\n", c.fullUse())
		writeln(w)
		writeln(w, "Available Commands:")

		type entry struct {
			name  string
			short string
		}
		entries := []entry{{name: "help", short: "Help about any command"}}
		for _, child := range c.children {
			entries = append(entries, entry{name: child.commandName(), short: child.Short})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].name < entries[j].name })

		width := 0
		for _, e := range entries {
			if len(e.name) > width {
				width = len(e.name)
			}
		}
		for _, e := range entries {
			writef(w, "  %-*s  %s\n", width, e.name, e.short)
		}
		writeln(w)
		return
	}

	writef(w, "  %s [flags]\n", c.fullUse())
	writeln(w)
}

func (c *Command) printFlags(w io.Writer) {
	all := c.root().PersistentFlags().sortedMetas()

	if c.parent == nil {
		writeln(w, "Flags:")
		c.writeFlagLines(w, all)
		return
	}

	var local, global []*flagMeta
	for _, meta := range all {
		if meta.name == "help" {
			local = append(local, meta)
		} else {
			global = append(global, meta)
		}
	}

	writeln(w, "Flags:")
	c.writeFlagLines(w, local)
	if len(global) > 0 {
		writeln(w)
		writeln(w, "Global Flags:")
		c.writeFlagLines(w, global)
	}
}

func (c *Command) writeFlagLines(w io.Writer, metas []*flagMeta) {
	if len(metas) == 0 {
		return
	}

	lines := make([]string, 0, len(metas))
	usages := make([]string, 0, len(metas))
	width := 0

	for _, meta := range metas {
		line := "      --" + meta.name
		if meta.shorthand != "" {
			line = "  -" + meta.shorthand + ", --" + meta.name
		}
		if kind := meta.typeName(); kind != "" {
			line += " " + kind
		}
		if len(line) > width {
			width = len(line)
		}

		usage := meta.usage
		if meta.name == "help" {
			usage = "help for " + c.commandName()
		}
		if suffix := meta.defaultSuffix(); suffix != "" {
			usage += suffix
		}
		lines = append(lines, line)
		usages = append(usages, usage)
	}

	for i, line := range lines {
		writef(w, "%-*s   %s\n", width, line, usages[i])
	}
}

// ==================== FlagSet ====================

type flagKind int

const (
	flagString flagKind = iota
	flagStringSlice
	flagBool
)

type flagMeta struct {
	name      string
	shorthand string
	usage     string
	kind      flagKind

	defString string
	defBool   bool

	stringValue *string
	sliceValue  *[]string
	boolValue   *bool
}

func (m *flagMeta) typeName() string {
	switch m.kind {
	case flagString:
		return "string"
	case flagStringSlice:
		return "strings"
	default:
		return ""
	}
}

func (m *flagMeta) defaultSuffix() string {
	switch m.kind {
	case flagString:
		if m.defString != "" {
			return fmt.Sprintf(" (default %q)", m.defString)
		}
	case flagBool:
		if m.defBool {
			return " (default true)"
		}
	}
	return ""
}

// FlagSet 全局 flag 集合
type FlagSet struct {
	fs      *flag.FlagSet
	metas   []*flagMeta
	byName  map[string]*flagMeta
	byShort map[string]*flagMeta
}

func newFlagSet(name string) *FlagSet {
	f := &FlagSet{
		fs:      flag.NewFlagSet(name, flag.ContinueOnError),
		byName:  make(map[string]*flagMeta),
		byShort: make(map[string]*flagMeta),
	}
	f.fs.SetOutput(io.Discard)
	f.fs.Usage = func() {}
	f.BoolP("help", "h", false, "help for "+name)
	return f
}

func (f *FlagSet) register(meta *flagMeta) {
	f.metas = append(f.metas, meta)
	f.byName[meta.name] = meta
	if meta.shorthand != "" {
		f.byShort[meta.shorthand] = meta
	}

	switch meta.kind {
	case flagString:
		f.fs.StringVar(meta.stringValue, meta.name, meta.defString, meta.usage)
		if meta.shorthand != "" {
			f.fs.StringVar(meta.stringValue, meta.shorthand, meta.defString, meta.usage)
		}
	case flagStringSlice:
		value := &stringSliceValue{target: meta.sliceValue}
		f.fs.Var(value, meta.name, meta.usage)
		if meta.shorthand != "" {
			f.fs.Var(value, meta.shorthand, meta.usage)
		}
	case flagBool:
		f.fs.BoolVar(meta.boolValue, meta.name, meta.defBool, meta.usage)
		if meta.shorthand != "" {
			f.fs.BoolVar(meta.boolValue, meta.shorthand, meta.defBool, meta.usage)
		}
	}
}

// String 定义字符串 flag（无短选项）
func (f *FlagSet) String(name, value, usage string) *string {
	return f.StringP(name, "", value, usage)
}

// StringP 定义字符串 flag
func (f *FlagSet) StringP(name, shorthand, value, usage string) *string {
	target := new(string)
	*target = value
	f.register(&flagMeta{
		name: name, shorthand: shorthand, usage: usage, kind: flagString,
		defString: value, stringValue: target,
	})
	return target
}

// StringSliceP 定义可重复、逗号分隔的字符串列表 flag
func (f *FlagSet) StringSliceP(name, shorthand string, value []string, usage string) *[]string {
	target := new([]string)
	*target = append([]string(nil), value...)
	f.register(&flagMeta{
		name: name, shorthand: shorthand, usage: usage, kind: flagStringSlice,
		sliceValue: target,
	})
	return target
}

// Bool 定义布尔 flag
func (f *FlagSet) Bool(name string, value bool, usage string) *bool {
	return f.BoolP(name, "", value, usage)
}

// BoolP 定义带短选项的布尔 flag
func (f *FlagSet) BoolP(name, shorthand string, value bool, usage string) *bool {
	target := new(bool)
	*target = value
	f.register(&flagMeta{
		name: name, shorthand: shorthand, usage: usage, kind: flagBool,
		defBool: value, boolValue: target,
	})
	return target
}

// GetString 获取字符串 flag 的值
func (f *FlagSet) GetString(name string) (string, error) {
	meta := f.byName[name]
	if meta == nil || meta.kind != flagString {
		return "", fmt.Errorf("flag %q not found", name)
	}
	return *meta.stringValue, nil
}

// GetStringSlice 获取字符串列表 flag 的值
func (f *FlagSet) GetStringSlice(name string) ([]string, error) {
	meta := f.byName[name]
	if meta == nil || meta.kind != flagStringSlice {
		return nil, fmt.Errorf("flag %q not found", name)
	}
	return append([]string(nil), (*meta.sliceValue)...), nil
}

// GetBool 获取布尔 flag 的值
func (f *FlagSet) GetBool(name string) (bool, error) {
	meta := f.byName[name]
	if meta == nil || meta.kind != flagBool {
		return false, fmt.Errorf("flag %q not found", name)
	}
	return *meta.boolValue, nil
}

func (f *FlagSet) sortedMetas() []*flagMeta {
	result := append([]*flagMeta(nil), f.metas...)
	sort.Slice(result, func(i, j int) bool { return result[i].name < result[j].name })
	return result
}

// parse 解析参数，允许 flag 与位置参数混排
func (f *FlagSet) parse(args []string) ([]string, error) {
	var flagArgs, positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			positional = append(positional, args[i+1:]...)
			break
		}
		if arg == "-" || !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
			continue
		}

		if strings.HasPrefix(arg, "--") {
			body := arg[2:]
			name, value, hasValue := strings.Cut(body, "=")
			meta := f.byName[name]
			if meta == nil {
				return nil, fmt.Errorf("unknown flag: --%s", name)
			}

			if meta.kind == flagBool {
				if hasValue {
					flagArgs = append(flagArgs, "--"+name+"="+value)
				} else {
					flagArgs = append(flagArgs, "--"+name)
				}
				continue
			}

			if hasValue {
				flagArgs = append(flagArgs, "--"+name+"="+value)
				continue
			}
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag needs an argument: --%s", name)
			}
			i++
			flagArgs = append(flagArgs, "--"+name, args[i])
			continue
		}

		// 单短横线：短选项聚合，例如 -y、-yh、-o out、-o=out、-oout
		rest := arg[1:]
		for len(rest) > 0 {
			shorthand := rest[:1]
			meta := f.byShort[shorthand]
			if meta == nil {
				return nil, fmt.Errorf("unknown shorthand flag: %q in -%s", rune(shorthand[0]), arg[1:])
			}
			rest = rest[1:]

			if meta.kind == flagBool {
				if strings.HasPrefix(rest, "=") {
					flagArgs = append(flagArgs, "-"+shorthand+"="+rest[1:])
					rest = ""
					continue
				}
				flagArgs = append(flagArgs, "-"+shorthand)
				continue
			}

			if strings.HasPrefix(rest, "=") {
				flagArgs = append(flagArgs, "-"+shorthand+"="+rest[1:])
				rest = ""
				continue
			}
			if rest != "" {
				flagArgs = append(flagArgs, "-"+shorthand+"="+rest)
				rest = ""
				continue
			}
			if i+1 >= len(args) {
				return nil, fmt.Errorf("flag needs an argument: %q in -%s", shorthand, shorthand)
			}
			i++
			flagArgs = append(flagArgs, "-"+shorthand, args[i])
		}
	}

	if err := f.fs.Parse(flagArgs); err != nil {
		return nil, err
	}
	return positional, nil
}

// stringSliceValue 与 pflag StringSlice 行为一致：首次设置替换默认值，之后追加
type stringSliceValue struct {
	target  *[]string
	changed bool
}

func (s *stringSliceValue) String() string {
	if s.target == nil {
		return ""
	}
	return strings.Join(*s.target, ",")
}

func (s *stringSliceValue) Set(value string) error {
	values, err := readAsCSV(value)
	if err != nil {
		return err
	}
	if !s.changed {
		*s.target = values
		s.changed = true
	} else {
		*s.target = append(*s.target, values...)
	}
	return nil
}

func readAsCSV(value string) ([]string, error) {
	if value == "" {
		return []string{}, nil
	}
	return csv.NewReader(strings.NewReader(value)).Read()
}
