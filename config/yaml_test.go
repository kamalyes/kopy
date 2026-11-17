/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\config\yaml_test.go
 * @Description: YAML 子集解析器与配置加载测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const fullTemplateYAML = `# 模板配置
name: demo-template
description: "模板描述 | Demo # not a comment"
version: 1.0.0
author: kamalyes
tags: [go, grpc, "cli, tool"]

delimiters:
  - "<%"
  - "%>"

questions:
  - name: ProjectName
    prompt: "项目名称 | Project name"
    type: string
    default: MyProject
    required: true

  - name: ServiceName
    prompt: "服务名称 | Service name"
    type: string
    default: "<% .ProjectName | snake %>-service"

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

ignore:
  - "kopy.yaml"
  - ".git/**"
  - '<% if ne .UseRedis "true" %>cache/redis.go<% end %>'

rename:
  - pattern: "__service_name__"
    value: "<% .ServiceName | snake %>"

hooks:
  after:
    - echo "hook ran"
`

func TestParseTemplateConfig(t *testing.T) {
	root, err := parseYAML([]byte(fullTemplateYAML))
	if err != nil {
		t.Fatalf("parseYAML() error = %v", err)
	}

	cfg, err := decodeTemplateConfig(root)
	if err != nil {
		t.Fatalf("decodeTemplateConfig() error = %v", err)
	}

	if cfg.Name != "demo-template" {
		t.Errorf("Name = %q, want %q", cfg.Name, "demo-template")
	}
	if cfg.Description != "模板描述 | Demo # not a comment" {
		t.Errorf("Description = %q", cfg.Description)
	}
	if cfg.Version != "1.0.0" {
		t.Errorf("Version = %q, want %q", cfg.Version, "1.0.0")
	}
	if cfg.Author != "kamalyes" {
		t.Errorf("Author = %q, want %q", cfg.Author, "kamalyes")
	}
	if got, want := strings.Join(cfg.Tags, ","), "go,grpc,cli, tool"; got != want {
		t.Errorf("Tags = %q, want %q", got, want)
	}
	if got, want := strings.Join(cfg.Delimiters, ","), "<%,%>"; got != want {
		t.Errorf("Delimiters = %q, want %q", got, want)
	}
	if len(cfg.Questions) != 5 {
		t.Fatalf("len(Questions) = %d, want 5", len(cfg.Questions))
	}

	first := cfg.Questions[0]
	if first.Name != "ProjectName" || first.Type != QuestionTypeString || first.Default != "MyProject" || !first.Required {
		t.Errorf("Questions[0] = %+v", first)
	}

	selectQ := cfg.Questions[2]
	if selectQ.Type != QuestionTypeSelect {
		t.Errorf("Questions[2].Type = %q, want select", selectQ.Type)
	}
	if got, want := strings.Join(selectQ.Choices, ","), "1.23,1.24,1.25"; got != want {
		t.Errorf("Questions[2].Choices = %q, want %q", got, want)
	}

	if len(cfg.Ignore) != 3 {
		t.Fatalf("len(Ignore) = %d, want 3", len(cfg.Ignore))
	}
	if cfg.Ignore[2] != `<% if ne .UseRedis "true" %>cache/redis.go<% end %>` {
		t.Errorf("Ignore[2] = %q", cfg.Ignore[2])
	}

	if len(cfg.Rename) != 1 || cfg.Rename[0].Pattern != "__service_name__" {
		t.Errorf("Rename = %+v", cfg.Rename)
	}

	if got, want := strings.Join(cfg.Hooks.After, ","), `echo "hook ran"`; got != want {
		t.Errorf("Hooks.After = %q, want %q", got, want)
	}
}

func TestParseYAMLFeatures(t *testing.T) {
	input := `---
name: features
empty_list: []
flow_map: {after: [echo one, echo two]}
literal: |
  line one
  # literal hash
  line three
folded: >
  folded one
  folded two

  new paragraph
null_value: null
tilde_value: ~
plain_hash: value#kept
plain_url: https://example.com/a#b
quoted: 'it''s ok'
escaped: "tab\tnewline\nunicode\u4F60"
`
	root, err := parseYAML([]byte(input))
	if err != nil {
		t.Fatalf("parseYAML() error = %v", err)
	}

	m, ok := root.(map[string]any)
	if !ok {
		t.Fatalf("root type = %T, want map[string]any", root)
	}

	if got, want := yamlString(m["literal"]), "line one\n# literal hash\nline three\n"; got != want {
		t.Errorf("literal = %q, want %q", got, want)
	}
	if got, want := yamlString(m["folded"]), "folded one folded two\nnew paragraph\n"; got != want {
		t.Errorf("folded = %q, want %q", got, want)
	}
	if m["null_value"] != nil {
		t.Errorf("null_value = %#v, want nil", m["null_value"])
	}
	if m["tilde_value"] != nil {
		t.Errorf("tilde_value = %#v, want nil", m["tilde_value"])
	}
	if got, want := yamlString(m["plain_hash"]), "value#kept"; got != want {
		t.Errorf("plain_hash = %q, want %q", got, want)
	}
	if got, want := yamlString(m["plain_url"]), "https://example.com/a#b"; got != want {
		t.Errorf("plain_url = %q, want %q", got, want)
	}
	if got, want := yamlString(m["quoted"]), "it's ok"; got != want {
		t.Errorf("quoted = %q, want %q", got, want)
	}
	if got, want := yamlString(m["escaped"]), "tab\tnewline\nunicode你"; got != want {
		t.Errorf("escaped = %q, want %q", got, want)
	}

	flowMap, ok := m["flow_map"].(map[string]any)
	if !ok {
		t.Fatalf("flow_map type = %T", m["flow_map"])
	}
	after, err := yamlStringList(flowMap["after"])
	if err != nil {
		t.Fatalf("flow_map.after: %v", err)
	}
	if got, want := strings.Join(after, ","), "echo one,echo two"; got != want {
		t.Errorf("flow_map.after = %q, want %q", got, want)
	}

	empty, err := yamlStringList(m["empty_list"])
	if err != nil {
		t.Fatalf("empty_list: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("empty_list = %#v, want empty", empty)
	}
}

func TestParseYAMLErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "tab indentation",
			input:   "name: x\n\tquestions: []\n",
			wantErr: "tab character used for indentation",
		},
		{
			name:    "duplicate key",
			input:   "name: x\nname: y\n",
			wantErr: "already defined",
		},
		{
			name:    "unterminated quote",
			input:   "name: \"x\n",
			wantErr: "unterminated quoted string",
		},
		{
			name:    "multiple documents",
			input:   "name: x\n---\nname: y\n",
			wantErr: "multiple YAML documents",
		},
		{
			name:    "content after document end",
			input:   "name: x\n...\nname: y\n",
			wantErr: "multiple YAML documents",
		},
		{
			name:    "unexpected indentation",
			input:   "name: x\n  description: y\n",
			wantErr: "unexpected indentation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseYAML([]byte(tt.input))
			if err == nil {
				t.Fatal("parseYAML() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestParseYAMLIndentlessSequence(t *testing.T) {
	input := `questions:
- name: ServiceName
  type: select
  choices:
  - a
  - b
ignore:
- one
- two
`
	root, err := parseYAML([]byte(input))
	if err != nil {
		t.Fatalf("parseYAML() error = %v", err)
	}
	cfg, err := decodeTemplateConfig(root)
	if err != nil {
		t.Fatalf("decodeTemplateConfig() error = %v", err)
	}
	if len(cfg.Questions) != 1 {
		t.Fatalf("len(Questions) = %d, want 1", len(cfg.Questions))
	}
	if got, want := strings.Join(cfg.Questions[0].Choices, ","), "a,b"; got != want {
		t.Errorf("Choices = %q, want %q", got, want)
	}
	if got, want := strings.Join(cfg.Ignore, ","), "one,two"; got != want {
		t.Errorf("Ignore = %q, want %q", got, want)
	}
}

func TestParseYAMLBlockScalarSequenceItem(t *testing.T) {
	input := "hooks:\n  after:\n    - |\n      echo one\n      echo two\n"
	root, err := parseYAML([]byte(input))
	if err != nil {
		t.Fatalf("parseYAML() error = %v", err)
	}
	cfg, err := decodeTemplateConfig(root)
	if err != nil {
		t.Fatalf("decodeTemplateConfig() error = %v", err)
	}
	if got, want := strings.Join(cfg.Hooks.After, "|"), "echo one\necho two\n"; got != want {
		t.Errorf("Hooks.After = %q, want %q", got, want)
	}
}

func TestParseYAMLTopLevelFlowMapping(t *testing.T) {
	root, err := parseYAML([]byte("{name: flow, tags: [a, b]}\n"))
	if err != nil {
		t.Fatalf("parseYAML() error = %v", err)
	}
	cfg, err := decodeTemplateConfig(root)
	if err != nil {
		t.Fatalf("decodeTemplateConfig() error = %v", err)
	}
	if cfg.Name != "flow" {
		t.Errorf("Name = %q, want flow", cfg.Name)
	}
	if got, want := strings.Join(cfg.Tags, ","), "a,b"; got != want {
		t.Errorf("Tags = %q, want %q", got, want)
	}
}

func TestParseYAMLMoreEdges(t *testing.T) {
	input := "\ufeff---\r\n" +
		`"name": quoted-key` + "\r\n" +
		"literal_strip: |-\n  a\n  b\n" +
		"folded_strip: >-\n  a\n  b\n" +
		"empty_map: {}\n" +
		"...\n"

	root, err := parseYAML([]byte(input))
	if err != nil {
		t.Fatalf("parseYAML() error = %v", err)
	}
	m := root.(map[string]any)

	if got, want := yamlString(m["name"]), "quoted-key"; got != want {
		t.Errorf("name = %q, want %q", got, want)
	}
	if got, want := yamlString(m["literal_strip"]), "a\nb"; got != want {
		t.Errorf("literal_strip = %q, want %q", got, want)
	}
	if got, want := yamlString(m["folded_strip"]), "a b"; got != want {
		t.Errorf("folded_strip = %q, want %q", got, want)
	}
	if empty, ok := m["empty_map"].(map[string]any); !ok || len(empty) != 0 {
		t.Errorf("empty_map = %#v, want empty map", m["empty_map"])
	}
}

func TestDecodeTemplateConfigErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name:    "top level not mapping",
			input:   "- a\n- b\n",
			wantErr: "top-level YAML value must be a mapping",
		},
		{
			name:    "questions not list",
			input:   "name: x\nquestions: nope\n",
			wantErr: "questions must be a list",
		},
		{
			name:    "question not mapping",
			input:   "name: x\nquestions:\n  - nope\n",
			wantErr: "questions[0] must be a mapping",
		},
		{
			name:    "choices not list",
			input:   "name: x\nquestions:\n  - name: a\n    choices: nope\n",
			wantErr: "questions[0].choices",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := parseYAML([]byte(tt.input))
			if err != nil {
				t.Fatalf("parseYAML() error = %v", err)
			}
			_, err = decodeTemplateConfig(root)
			if err == nil {
				t.Fatal("decodeTemplateConfig() error = nil, want error")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestLoadFromDir(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, ConfigFileName), []byte("name: demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadFromDir(dir)
	if err != nil {
		t.Fatalf("LoadFromDir() error = %v", err)
	}
	if got, want := strings.Join(cfg.Delimiters, ","), "[[,]]"; got != want {
		t.Errorf("default Delimiters = %q, want %q", got, want)
	}
	if got, want := strings.Join(cfg.Ignore, ","), "kopy.yaml,.git/**"; got != want {
		t.Errorf("default Ignore = %q, want %q", got, want)
	}
}

func TestLoadFromDirErrors(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		_, err := LoadFromDir(t.TempDir())
		if err == nil || !strings.Contains(err.Error(), "读取配置文件失败") {
			t.Fatalf("error = %v, want read failure", err)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ConfigFileName), []byte("description: x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := LoadFromDir(dir)
		if err == nil || !strings.Contains(err.Error(), "模板名称为必填项") {
			t.Fatalf("error = %v, want name required", err)
		}
	})

	t.Run("parse error includes reason", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, ConfigFileName), []byte("name: \"x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := LoadFromDir(dir)
		if err == nil || !strings.Contains(err.Error(), "unterminated quoted string") {
			t.Fatalf("error = %v, want parse reason", err)
		}
	})
}
