/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\errors\errors_test.go
 * @Description: 类型化错误测试
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package errors

import (
	stderrors "errors"
	"strings"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewConfigQuestionDuplicateNameError("2", "Name")
	if got, want := err.Error(), `问题[2]: 名称重复 "Name"`; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}

	var base BaseError
	if !stderrors.As(err, &base) {
		t.Fatal("errors.As() = false, want BaseError")
	}
	if base.GetType() != ErrTypeConfigQuestionDuplicateName {
		t.Errorf("GetType() = %d, want %d", base.GetType(), ErrTypeConfigQuestionDuplicateName)
	}
}

func TestNewErrorUnknownType(t *testing.T) {
	err := NewError(ErrorType(-1))
	if got, want := err.Error(), "unknown error"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestWrapError(t *testing.T) {
	if err := WrapError("context"); err == nil || err.Error() != "context" {
		t.Errorf("WrapError without cause = %v, want context", err)
	}

	cause := stderrors.New("cause")
	err := WrapError("context", cause)
	if !stderrors.Is(err, cause) {
		t.Error("errors.Is() = false, want wrapped cause")
	}
	if !strings.HasPrefix(err.Error(), "context: cause") {
		t.Errorf("Error() = %q", err)
	}
}
