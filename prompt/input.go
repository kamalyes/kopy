/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\prompt\input.go
 * @Description: 各类型输入方法 - string/bool/int/select
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package prompt

import (
	"fmt"
	"strconv"
	"strings"
)

// askString 文本输入
func (p *Prompter) askString(prompt, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("  ❓ %s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("  ❓ %s: ", prompt)
	}

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal, nil
	}
	return input, nil
}

// askBool 布尔选择
func (p *Prompter) askBool(prompt, defaultVal string) (string, error) {
	hint := "y/N"
	if defaultVal == "true" || defaultVal == "y" || defaultVal == "yes" {
		hint = "Y/n"
	}

	fmt.Printf("  ❓ %s [%s]: ", prompt, hint)
	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		if defaultVal == "true" || defaultVal == "y" || defaultVal == "yes" {
			return "true", nil
		}
		return "false", nil
	}

	if input == "y" || input == "yes" || input == "true" {
		return "true", nil
	}
	return "false", nil
}

// askInt 整数输入
func (p *Prompter) askInt(prompt, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("  ❓ %s [%s]: ", prompt, defaultVal)
	} else {
		fmt.Printf("  ❓ %s: ", prompt)
	}

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal, nil
	}

	if _, err := strconv.Atoi(input); err != nil {
		fmt.Println("  ⚠ Please enter a valid integer")
		return p.askInt(prompt, defaultVal)
	}
	return input, nil
}

// askSelect 下拉选择
func (p *Prompter) askSelect(prompt string, choices []string, defaultVal string) (string, error) {
	fmt.Printf("  ❓ %s:\n", prompt)
	for i, choice := range choices {
		marker := "  "
		if choice == defaultVal {
			marker = "→ "
		}
		fmt.Printf("    %s %d) %s\n", marker, i+1, choice)
	}
	fmt.Printf("  Enter number [%s]: ", defaultVal)

	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal, nil
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(choices) {
		fmt.Println("  ⚠ Invalid choice")
		return p.askSelect(prompt, choices, defaultVal)
	}
	return choices[idx-1], nil
}
