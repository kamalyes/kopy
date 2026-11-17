/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\bootstrap\cli_test.go
 * @Description: 轻量命令与参数解析测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package bootstrap

import (
	"bytes"
	"strings"
	"testing"
)

type testRun struct {
	called    bool
	args      []string
	output    string
	vars      []string
	dryRun    bool
	overwrite bool
	yes       bool
	hooks     string
	out       bytes.Buffer
	err       bytes.Buffer
}

func newTestCommand(run *testRun) *Command {
	root := &Command{Use: "kopy"}
	flags := root.PersistentFlags()
	flags.StringP("output", "o", ".", "output")
	flags.StringSliceP("var", "v", nil, "vars")
	flags.Bool("dry-run", false, "dry run")
	flags.Bool("overwrite", false, "overwrite")
	flags.BoolP("yes", "y", false, "yes")
	flags.String("hooks", "ask", "hooks")

	child := &Command{
		Use:   "new <src>",
		Short: "new",
		Args:  ExactArgs(1),
		Run: func(cmd *Command, args []string) {
			run.called = true
			run.args = args
			run.output, _ = cmd.Flags().GetString("output")
			run.vars, _ = cmd.Flags().GetStringSlice("var")
			run.dryRun, _ = cmd.Flags().GetBool("dry-run")
			run.overwrite, _ = cmd.Flags().GetBool("overwrite")
			run.yes, _ = cmd.Flags().GetBool("yes")
			run.hooks, _ = cmd.Flags().GetString("hooks")
		},
	}
	root.AddCommand(child)
	root.Out = &run.out
	root.Err = &run.err
	return root
}

func TestCLIParsesInterspersedFlags(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	err := root.ExecuteArgs([]string{"new", "src", "-o", "out", "-y", "--var", "a=b", "-v", "c=d", "--dry-run"})
	if err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if !run.called {
		t.Fatal("Run was not called")
	}
	if got, want := strings.Join(run.args, ","), "src"; got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
	if run.output != "out" {
		t.Errorf("output = %q, want out", run.output)
	}
	if !run.yes || !run.dryRun {
		t.Errorf("yes = %v, dryRun = %v, want true/true", run.yes, run.dryRun)
	}
	if got, want := strings.Join(run.vars, ","), "a=b,c=d"; got != want {
		t.Errorf("vars = %q, want %q", got, want)
	}
}

func TestCLIParsesFlagsBeforeSubcommand(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"-o", "out", "--hooks", "never", "new", "src"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if run.output != "out" || run.hooks != "never" {
		t.Errorf("output = %q, hooks = %q", run.output, run.hooks)
	}
}

func TestCLICommaSeparatedVars(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"new", "src", "--var=a=b,c=d"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if got, want := strings.Join(run.vars, ","), "a=b,c=d"; got != want {
		t.Errorf("vars = %q, want %q", got, want)
	}
}

func TestCLIShorthandAttachedValue(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"new", "src", "-vA=B", "-oout"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if got, want := strings.Join(run.vars, ","), "A=B"; got != want {
		t.Errorf("vars = %q, want %q", got, want)
	}
	if run.output != "out" {
		t.Errorf("output = %q, want out", run.output)
	}
}

func TestCLIShorthandEmptyValue(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"new", "src", "-o=", "-v="}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if !run.called {
		t.Fatal("Run was not called")
	}
	if got, want := strings.Join(run.args, ","), "src"; got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
	if run.output != "" {
		t.Errorf("output = %q, want empty", run.output)
	}
	if len(run.vars) != 0 {
		t.Errorf("vars = %#v, want empty", run.vars)
	}
}

func TestCLIBoolExplicitFalse(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"new", "src", "--yes=false", "-y=false", "--overwrite=true"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if run.yes {
		t.Error("yes = true, want false")
	}
	if !run.overwrite {
		t.Error("overwrite = false, want true")
	}
}

func TestCLIHelp(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"--help"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if run.called {
		t.Error("Run was called on help")
	}
	out := run.out.String()
	for _, want := range []string{"Usage:", "Available Commands:", "new", "Flags:"} {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing %q:\n%s", want, out)
		}
	}
}

func TestCLIChildHelp(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"new", "--help"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if run.called {
		t.Error("Run was called on help")
	}
	out := run.out.String()
	for _, want := range []string{"kopy new <src> [flags]", "Global Flags:", "help for new"} {
		if !strings.Contains(out, want) {
			t.Errorf("child help missing %q:\n%s", want, out)
		}
	}
}

func TestCLIHelpCommand(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"help", "new"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if !strings.Contains(run.out.String(), "kopy new <src> [flags]") {
		t.Errorf("help new output:\n%s", run.out.String())
	}
}

func TestCLICombinedShorthands(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"new", "src", "-yh"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if run.called {
		t.Error("Run was called with -yh")
	}
	if !strings.Contains(run.out.String(), "kopy new <src> [flags]") {
		t.Errorf("output:\n%s", run.out.String())
	}
}

func TestCLIDoubleDash(t *testing.T) {
	run := &testRun{}
	root := newTestCommand(run)

	if err := root.ExecuteArgs([]string{"new", "--", "-weird"}); err != nil {
		t.Fatalf("ExecuteArgs() error = %v", err)
	}
	if got, want := strings.Join(run.args, ","), "-weird"; got != want {
		t.Errorf("args = %q, want %q", got, want)
	}
}

func TestCLIErrors(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantErr    string
		wantStderr string
	}{
		{
			name:       "missing arg",
			args:       []string{"new"},
			wantErr:    "accepts 1 arg(s), received 0",
			wantStderr: "Usage:",
		},
		{
			name:       "too many args",
			args:       []string{"new", "a", "b"},
			wantErr:    "accepts 1 arg(s), received 2",
			wantStderr: "Usage:",
		},
		{
			name:       "unknown command",
			args:       []string{"nope"},
			wantErr:    `unknown command "nope"`,
			wantStderr: "Run 'kopy --help' for usage.",
		},
		{
			name:       "unknown flag",
			args:       []string{"--bad"},
			wantErr:    "unknown flag: --bad",
			wantStderr: "Usage:",
		},
		{
			name:       "unknown shorthand",
			args:       []string{"-z"},
			wantErr:    "unknown shorthand flag",
			wantStderr: "Usage:",
		},
		{
			name:       "missing flag value",
			args:       []string{"--output"},
			wantErr:    "flag needs an argument",
			wantStderr: "Usage:",
		},
		{
			name:       "unknown help topic",
			args:       []string{"help", "nope"},
			wantErr:    `unknown help topic "nope"`,
			wantStderr: "Run 'kopy --help' for usage.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run := &testRun{}
			root := newTestCommand(run)

			err := root.ExecuteArgs(tt.args)
			if err == nil {
				t.Fatal("ExecuteArgs() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want containing %q", err, tt.wantErr)
			}
			if !strings.Contains(run.err.String(), tt.wantStderr) {
				t.Errorf("stderr = %q, want containing %q", run.err.String(), tt.wantStderr)
			}
		})
	}
}
