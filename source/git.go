/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\source\git.go
 * @Description: Git 仓库模板源 - 克隆仓库并查找模板目录
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package source

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	kopyerrors "github.com/kamalyes/kopy/errors"
	kopylog "github.com/kamalyes/kopy/log"
)

type gitSource struct {
	url    string
	tmpDir string
}

func (s *gitSource) Resolve() (string, error) {
	tmpDir, err := os.MkdirTemp("", "kopy-*")
	if err != nil {
		return "", kopyerrors.NewSourceTempDirFailedError(err.Error())
	}
	s.tmpDir = tmpDir

	url, branch := parseGitURL(s.url)

	args := []string{"clone", "--depth", "1"}
	if branch != "" {
		args = append(args, "--branch", branch)
	}
	args = append(args, url, tmpDir)

	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		kopylog.Error("git clone %s failed: %s", url, string(output))
		return "", kopyerrors.NewSourceGitCloneFailedError(url)
	}

	return findTemplateDir(tmpDir), nil
}

func (s *gitSource) Cleanup() {
	if s.tmpDir != "" {
		os.RemoveAll(s.tmpDir)
	}
}

func parseGitURL(raw string) (url, branch string) {
	parts := strings.SplitN(raw, "#", 2)
	url = parts[0]
	if len(parts) == 2 {
		branch = parts[1]
	}
	return
}

func findTemplateDir(root string) string {
	configPath := filepath.Join(root, "kopy.yaml")
	if _, err := os.Stat(configPath); err == nil {
		return root
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return root
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		subConfig := filepath.Join(root, entry.Name(), "kopy.yaml")
		if _, err := os.Stat(subConfig); err == nil {
			return filepath.Join(root, entry.Name())
		}
	}

	return root
}
