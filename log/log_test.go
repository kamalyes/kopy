/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\log\log_test.go
 * @Description: 日志输出测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package log

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestLogLevels(t *testing.T) {
	var buffer bytes.Buffer
	SetOutput(&buffer)
	defer SetOutput(os.Stdout)

	Debug("debug %s", "message")
	Info("info %s", "message")
	Warn("warn %s", "message")
	Error("error %s", "message")
	Fatal("fatal %s", "message")

	output := buffer.String()
	for _, want := range []string{
		"[KOPY]",
		"🐛 [DEBUG]",
		"debug message",
		"ℹ️ [INFO]",
		"info message",
		"⚠️ [WARN]",
		"warn message",
		"❌ [ERROR]",
		"error message",
		"💀 [FATAL]",
		"fatal message",
	} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestSetLogger(t *testing.T) {
	var buffer bytes.Buffer
	SetOutput(&buffer)
	defer SetOutput(os.Stdout)

	SetLogger(nil)
	if L() == nil {
		t.Fatal("L() = nil after SetLogger(nil)")
	}
}
