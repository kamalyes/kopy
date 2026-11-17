/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\engine\naming_test.go
 * @Description: 命名转换测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package engine

import "testing"

func TestTitleCase(t *testing.T) {
	tests := map[string]string{
		"hello world":   "Hello World",
		"hello-world":   "Hello-World",
		"hello_world":   "Hello_world",
		"hello.world":   "Hello.World",
		"":              "",
		"already Title": "Already Title",
	}
	for input, want := range tests {
		if got := titleCase(input); got != want {
			t.Errorf("titleCase(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNamingConversions(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		snake  string
		camel  string
		pascal string
		kebab  string
	}{
		{name: "pascal", input: "MyService", snake: "my_service", camel: "myService", pascal: "MyService", kebab: "my-service"},
		{name: "camel", input: "myService", snake: "my_service", camel: "myService", pascal: "MyService", kebab: "my-service"},
		{name: "snake", input: "my_service", snake: "my_service", camel: "myService", pascal: "MyService", kebab: "my-service"},
		{name: "kebab", input: "my-service", snake: "my_service", camel: "myService", pascal: "MyService", kebab: "my-service"},
		{name: "spaces", input: "my service", snake: "my_service", camel: "myService", pascal: "MyService", kebab: "my-service"},
		{name: "acronym", input: "HTTPServer", snake: "httpserver", camel: "hTTPServer", pascal: "HTTPServer", kebab: "httpserver"},
		{name: "mixed separators", input: "HTTP_server", snake: "http_server", camel: "httpServer", pascal: "HttpServer", kebab: "http-server"},
		{name: "double underscore", input: "my__service", snake: "my__service", camel: "myService", pascal: "MyService", kebab: "my--service"},
		{name: "empty", input: "", snake: "", camel: "", pascal: "", kebab: ""},
		{name: "trim", input: "  my_service  ", snake: "my_service", camel: "myService", pascal: "MyService", kebab: "my-service"},
		{name: "unicode", input: "用户_name", snake: "用户_name", camel: "用户Name", pascal: "用户Name", kebab: "用户-name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toSnakeCase(tt.input); got != tt.snake {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, got, tt.snake)
			}
			if got := toCamelCase(tt.input); got != tt.camel {
				t.Errorf("toCamelCase(%q) = %q, want %q", tt.input, got, tt.camel)
			}
			if got := toPascalCase(tt.input); got != tt.pascal {
				t.Errorf("toPascalCase(%q) = %q, want %q", tt.input, got, tt.pascal)
			}
			if got := toKebabCase(tt.input); got != tt.kebab {
				t.Errorf("toKebabCase(%q) = %q, want %q", tt.input, got, tt.kebab)
			}
		})
	}
}
