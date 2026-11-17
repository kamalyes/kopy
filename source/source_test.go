/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\source\source_test.go
 * @Description: 模板源解析测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package source

import "testing"

func TestNewShorthand(t *testing.T) {
	tests := []struct {
		raw    string
		url    string
		branch string
	}{
		{"@linmi-tech/kronos-service-template", "https://github.com/linmi-tech/kronos-service-template", ""},
		{"@linmi-tech/kronos-service-template@v1.0", "https://github.com/linmi-tech/kronos-service-template", "v1.0"},
		{"@linmi-tech/kronos-service-template@feature/foo", "https://github.com/linmi-tech/kronos-service-template", "feature/foo"},
		{"@linmi-tech/kronos-service-template#main", "https://github.com/linmi-tech/kronos-service-template", "main"},
	}

	for _, tt := range tests {
		src := New(tt.raw)
		gitSrc, ok := src.(*gitSource)
		if !ok {
			t.Errorf("New(%q) = %T, want *gitSource", tt.raw, src)
			continue
		}
		url, branch := parseGitURL(gitSrc.url)
		if url != tt.url || branch != tt.branch {
			t.Errorf("New(%q) = (%q, %q), want (%q, %q)", tt.raw, url, branch, tt.url, tt.branch)
		}
	}
}

func TestNewLocalAndURL(t *testing.T) {
	if _, ok := New("./my-template").(*localSource); !ok {
		t.Error("New(./my-template) should be a local source")
	}
	if _, ok := New("/tmp/my-template").(*localSource); !ok {
		t.Error("New(/tmp/my-template) should be a local source")
	}
	if _, ok := New("https://github.com/linmi-tech/kopy.git").(*gitSource); !ok {
		t.Error("New(https URL) should be a git source")
	}
	if _, ok := New("git@github.com:linmi-tech/kopy.git").(*gitSource); !ok {
		t.Error("New(git@ URL) should be a git source")
	}
}

func TestParseShorthandRejectsNonRepo(t *testing.T) {
	for _, raw := range []string{"@owner", "@owner/repo/sub", "@@owner/repo", "owner/repo"} {
		if _, _, ok := parseShorthand(raw); ok {
			t.Errorf("parseShorthand(%q) = true, want false", raw)
		}
	}
}
