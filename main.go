/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 23:29:19
 * @FilePath: \kopy\main.go
 * @Description: kopy 通用项目模板生成工具入口
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package main

import (
	"os"

	"github.com/kamalyes/kopy/bootstrap"
)

func main() {
	if err := bootstrap.Execute(); err != nil {
		os.Exit(1)
	}
}
