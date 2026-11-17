/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2026-04-19 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-04-19 00:00:00
 * @FilePath: \kopy\config\yaml.go
 * @Description: 轻量 YAML 子集解析器 - 仅依赖标准库，覆盖 kopy.yaml 所需语法
 *
 * 支持：
 *   - 块级映射与序列（缩进使用空格）
 *   - flow 序列 [a, b] 与 flow 映射 {k: v}
 *   - 单引号/双引号字符串与常见转义
 *   - 行注释、空行、--- 文档起始标记
 *   - 块标量 | 与 >（含 -/+ 截断指示符）
 *
 * 不支持锚点、别名、标签、多文档等 kopy.yaml 不需要的高级语法
 *
 * Copyright (c) 2026 by kamalyes, All Rights Reserved.
 */

package config

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

type yamlLine struct {
	raw    string
	indent int
	text   string
	number int
}

type yamlParser struct {
	lines []yamlLine
	pos   int
}

// parseYAML 将 YAML 文本解析为 map[string]any / []any / string / nil 组成的树
func parseYAML(data []byte) (any, error) {
	p, err := newYAMLParser(data)
	if err != nil {
		return nil, err
	}

	p.skipIgnorable()
	if p.pos >= len(p.lines) {
		return nil, nil
	}
	if p.lines[p.pos].indent != 0 {
		return nil, fmt.Errorf("line %d: unexpected indentation", p.lines[p.pos].number)
	}

	v, err := p.parseNode(0)
	if err != nil {
		return nil, err
	}

	p.skipIgnorable()
	if p.pos < len(p.lines) {
		return nil, fmt.Errorf("line %d: unexpected content", p.lines[p.pos].number)
	}
	return v, nil
}

func newYAMLParser(data []byte) (*yamlParser, error) {
	raw := strings.TrimPrefix(string(data), "\ufeff")
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")

	rawLines := strings.Split(raw, "\n")
	lines := make([]yamlLine, 0, len(rawLines))
	documentStarted := false
	hasContent := false

	for i, line := range rawLines {
		number := i + 1
		indent := 0
		for indent < len(line) && line[indent] == ' ' {
			indent++
		}

		text := strings.TrimRight(line[indent:], " \t")
		blank := strings.TrimSpace(line) == ""

		if !blank {
			switch text {
			case "---":
				if documentStarted || hasContent {
					return nil, fmt.Errorf("line %d: multiple YAML documents are not supported", number)
				}
				documentStarted = true
				continue
			case "...":
				for j := i + 1; j < len(rawLines); j++ {
					rest := strings.TrimSpace(rawLines[j])
					if rest != "" && !strings.HasPrefix(rest, "#") {
						return nil, fmt.Errorf("line %d: multiple YAML documents are not supported", j+1)
					}
				}
				return &yamlParser{lines: lines}, nil
			}

			if text[0] != '#' {
				text = strings.TrimRight(stripYAMLComment(text), " \t")
			}
			if text == "" || text[0] == '#' {
				text = ""
			} else {
				hasContent = true
			}
		}

		lines = append(lines, yamlLine{
			raw:    line,
			indent: indent,
			text:   text,
			number: number,
		})
	}

	return &yamlParser{lines: lines}, nil
}

// stripYAMLComment 去掉引号外的行注释
func stripYAMLComment(s string) string {
	inSingle, inDouble, escaped := false, false, false

	for i := 0; i < len(s); i++ {
		c := s[i]
		if inDouble {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				inDouble = false
			}
			continue
		}
		if inSingle {
			if c == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
			}
			continue
		}

		switch c {
		case '"':
			inDouble = true
		case '\'':
			inSingle = true
		case '#':
			if i == 0 || s[i-1] == ' ' || s[i-1] == '\t' {
				return s[:i]
			}
		}
	}
	return s
}

func (p *yamlParser) skipIgnorable() {
	for p.pos < len(p.lines) && p.lines[p.pos].text == "" {
		p.pos++
	}
}

func (p *yamlParser) parseNode(indent int) (any, error) {
	p.skipIgnorable()
	if p.pos >= len(p.lines) {
		return nil, nil
	}

	line := p.lines[p.pos]
	if line.indent < indent {
		return nil, nil
	}
	if line.indent > indent {
		return nil, fmt.Errorf("line %d: unexpected indentation", line.number)
	}
	if err := checkYAMLIndent(line); err != nil {
		return nil, err
	}
	if line.text[0] == '{' || line.text[0] == '[' {
		p.pos++
		return p.parseInlineValue(line.text, line.number)
	}
	if isYAMLSeqItem(line.text) {
		return p.parseSequence(indent)
	}
	return p.parseMapping(indent)
}

func checkYAMLIndent(line yamlLine) error {
	if line.indent < len(line.raw) && line.raw[line.indent] == '\t' {
		return fmt.Errorf("line %d: tab character used for indentation", line.number)
	}
	return nil
}

func (p *yamlParser) parseMapping(indent int) (map[string]any, error) {
	result := make(map[string]any)

	for {
		p.skipIgnorable()
		if p.pos >= len(p.lines) {
			break
		}

		line := p.lines[p.pos]
		if line.indent < indent {
			break
		}
		if line.indent > indent {
			return nil, fmt.Errorf("line %d: unexpected indentation", line.number)
		}
		if err := checkYAMLIndent(line); err != nil {
			return nil, err
		}
		if isYAMLSeqItem(line.text) {
			break
		}

		key, rest, ok := splitYAMLKeyValue(line.text)
		if !ok {
			return nil, fmt.Errorf("line %d: expected key: value", line.number)
		}
		p.pos++

		value, err := p.parseEntryValue(indent, rest, line.number)
		if err != nil {
			return nil, err
		}
		key = decodeYAMLKey(key)
		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("line %d: mapping key %q already defined", line.number, key)
		}
		result[key] = value
	}

	return result, nil
}

func (p *yamlParser) parseSequence(indent int) ([]any, error) {
	var result []any

	for {
		p.skipIgnorable()
		if p.pos >= len(p.lines) {
			break
		}

		line := p.lines[p.pos]
		if line.indent < indent {
			break
		}
		if line.indent > indent {
			return nil, fmt.Errorf("line %d: unexpected indentation", line.number)
		}
		if err := checkYAMLIndent(line); err != nil {
			return nil, err
		}
		if !isYAMLSeqItem(line.text) {
			break
		}

		rest := line.text[1:]
		contentIndent := line.indent + 1
		for contentIndent < len(line.raw) && line.raw[contentIndent] == ' ' {
			contentIndent++
		}
		rest = strings.TrimLeft(rest, " ")
		p.pos++

		var (
			item any
			err  error
		)

		switch {
		case rest == "":
			p.skipIgnorable()
			if p.pos < len(p.lines) && p.lines[p.pos].indent > indent {
				item, err = p.parseNode(p.lines[p.pos].indent)
			}
		case isYAMLSeqItem(rest):
			err = fmt.Errorf("line %d: nested inline sequences are not supported", line.number)
		case isYAMLBlockScalarHeader(rest):
			item, err = p.parseBlockScalar(indent, rest)
		default:
			if key, value, ok := splitYAMLKeyValue(rest); ok {
				item, err = p.parseSequenceMapping(contentIndent, key, value, line.number)
			} else {
				item, err = p.parseInlineValue(rest, line.number)
			}
		}
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

// parseSequenceMapping 解析 "- key: value" 起始的映射项
func (p *yamlParser) parseSequenceMapping(indent int, firstKey, firstValue string, lineNumber int) (map[string]any, error) {
	result := make(map[string]any)

	value, err := p.parseEntryValue(indent, firstValue, lineNumber)
	if err != nil {
		return nil, err
	}
	firstKey = decodeYAMLKey(firstKey)
	result[firstKey] = value

	for {
		p.skipIgnorable()
		if p.pos >= len(p.lines) {
			break
		}

		line := p.lines[p.pos]
		if line.indent < indent {
			break
		}
		if line.indent > indent {
			return nil, fmt.Errorf("line %d: unexpected indentation", line.number)
		}
		if err := checkYAMLIndent(line); err != nil {
			return nil, err
		}
		if isYAMLSeqItem(line.text) {
			break
		}

		key, rest, ok := splitYAMLKeyValue(line.text)
		if !ok {
			return nil, fmt.Errorf("line %d: expected key: value", line.number)
		}
		p.pos++

		value, err := p.parseEntryValue(indent, rest, line.number)
		if err != nil {
			return nil, err
		}
		key = decodeYAMLKey(key)
		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("line %d: mapping key %q already defined", line.number, key)
		}
		result[key] = value
	}

	return result, nil
}

func (p *yamlParser) parseEntryValue(parentIndent int, rest string, lineNumber int) (any, error) {
	rest = strings.TrimSpace(rest)

	if rest == "" {
		p.skipIgnorable()
		if p.pos < len(p.lines) {
			next := p.lines[p.pos]
			if next.indent > parentIndent {
				return p.parseNode(next.indent)
			}
			// YAML 允许无缩进序列：key:\n- item
			if next.indent == parentIndent && isYAMLSeqItem(next.text) {
				return p.parseSequence(parentIndent)
			}
		}
		return nil, nil
	}

	if isYAMLBlockScalarHeader(rest) {
		return p.parseBlockScalar(parentIndent, rest)
	}

	return p.parseInlineValue(rest, lineNumber)
}

func (p *yamlParser) parseInlineValue(s string, lineNumber int) (any, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}

	switch s[0] {
	case '[':
		value, end, err := parseYAMLFlowSequence(s, 0)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if strings.TrimSpace(s[end:]) != "" {
			return nil, fmt.Errorf("line %d: unexpected content after flow sequence", lineNumber)
		}
		return value, nil
	case '{':
		value, end, err := parseYAMLFlowMapping(s, 0)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if strings.TrimSpace(s[end:]) != "" {
			return nil, fmt.Errorf("line %d: unexpected content after flow mapping", lineNumber)
		}
		return value, nil
	case '"', '\'':
		value, end, err := parseYAMLQuoted(s, 0)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if strings.TrimSpace(s[end:]) != "" {
			return nil, fmt.Errorf("line %d: unexpected content after quoted string", lineNumber)
		}
		return value, nil
	default:
		if isYAMLNull(s) {
			return nil, nil
		}
		return s, nil
	}
}

func (p *yamlParser) parseBlockScalar(parentIndent int, header string) (string, error) {
	style := header[0]
	chomp := byte(0)
	explicitIndent := 0
	for i := 1; i < len(header); i++ {
		switch c := header[i]; {
		case c == '+' || c == '-':
			chomp = c
		case c >= '1' && c <= '9':
			explicitIndent = int(c - '0')
		}
	}

	contentIndent := -1
	if explicitIndent > 0 {
		contentIndent = parentIndent + explicitIndent
	}

	var content []string
	for p.pos < len(p.lines) {
		line := p.lines[p.pos]
		if strings.TrimSpace(line.raw) == "" {
			content = append(content, "")
			p.pos++
			continue
		}
		if line.indent <= parentIndent {
			break
		}
		if contentIndent < 0 {
			contentIndent = line.indent
		}
		if line.indent < contentIndent {
			return "", fmt.Errorf("line %d: invalid indentation in block scalar", line.number)
		}
		content = append(content, line.raw[contentIndent:])
		p.pos++
	}

	if contentIndent < 0 {
		return "", nil
	}

	var value string
	if style == '|' {
		value = strings.Join(content, "\n")
	} else {
		var builder strings.Builder
		previousBlank := false
		for i, line := range content {
			if line == "" {
				builder.WriteByte('\n')
				previousBlank = true
				continue
			}
			if i > 0 && !previousBlank {
				builder.WriteByte(' ')
			}
			builder.WriteString(line)
			previousBlank = false
		}
		value = builder.String()
	}

	value = strings.TrimRight(value, "\n")
	switch chomp {
	case '-':
		return value, nil
	case '+':
		trailing := 0
		for i := len(content) - 1; i >= 0 && content[i] == ""; i-- {
			trailing++
		}
		return value + "\n" + strings.Repeat("\n", trailing), nil
	default:
		return value + "\n", nil
	}
}

func isYAMLBlockScalarHeader(s string) bool {
	if s == "" || (s[0] != '|' && s[0] != '>') {
		return false
	}
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '+', '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return false
		}
	}
	return true
}

func isYAMLSeqItem(s string) bool {
	return s == "-" || strings.HasPrefix(s, "- ")
}

func isYAMLNull(s string) bool {
	switch s {
	case "null", "Null", "NULL", "~":
		return true
	}
	return false
}

func splitYAMLKeyValue(s string) (key, value string, ok bool) {
	inSingle, inDouble, escaped := false, false, false
	depth := 0

	for i := 0; i < len(s); i++ {
		c := s[i]
		if inDouble {
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == '"' {
				inDouble = false
			}
			continue
		}
		if inSingle {
			if c == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					i++
					continue
				}
				inSingle = false
			}
			continue
		}

		switch c {
		case '"':
			inDouble = true
		case '\'':
			inSingle = true
		case '[', '{':
			depth++
		case ']', '}':
			if depth > 0 {
				depth--
			}
		case ':':
			if depth == 0 && (i+1 == len(s) || s[i+1] == ' ' || s[i+1] == '\t') {
				key = strings.TrimSpace(s[:i])
				value = strings.TrimSpace(s[i+1:])
				return key, value, key != ""
			}
		}
	}
	return "", "", false
}

func decodeYAMLKey(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && (s[0] == '"' || s[0] == '\'') {
		if value, end, err := parseYAMLQuoted(s, 0); err == nil && strings.TrimSpace(s[end:]) == "" {
			return value
		}
	}
	return s
}

func parseYAMLFlowSequence(s string, start int) ([]any, int, error) {
	if start >= len(s) || s[start] != '[' {
		return nil, 0, fmt.Errorf("expected '['")
	}

	i := skipYAMLFlowSpace(s, start+1)
	items := make([]any, 0)
	if i < len(s) && s[i] == ']' {
		return items, i + 1, nil
	}

	for {
		value, next, err := parseYAMLFlowValue(s, i)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, value)
		i = skipYAMLFlowSpace(s, next)

		if i >= len(s) {
			return nil, 0, fmt.Errorf("unterminated flow sequence")
		}
		switch s[i] {
		case ',':
			i = skipYAMLFlowSpace(s, i+1)
			if i < len(s) && s[i] == ']' {
				return items, i + 1, nil
			}
		case ']':
			return items, i + 1, nil
		default:
			return nil, 0, fmt.Errorf("expected ',' or ']' in flow sequence")
		}
	}
}

func parseYAMLFlowMapping(s string, start int) (map[string]any, int, error) {
	if start >= len(s) || s[start] != '{' {
		return nil, 0, fmt.Errorf("expected '{'")
	}

	i := skipYAMLFlowSpace(s, start+1)
	result := make(map[string]any)
	if i < len(s) && s[i] == '}' {
		return result, i + 1, nil
	}

	for {
		var key string
		if i < len(s) && (s[i] == '"' || s[i] == '\'') {
			value, next, err := parseYAMLQuoted(s, i)
			if err != nil {
				return nil, 0, err
			}
			key = value
			i = skipYAMLFlowSpace(s, next)
		} else {
			keyStart := i
			for i < len(s) && s[i] != ':' {
				i++
			}
			key = strings.TrimSpace(s[keyStart:i])
		}

		if i >= len(s) || s[i] != ':' {
			return nil, 0, fmt.Errorf("expected ':' in flow mapping")
		}
		i = skipYAMLFlowSpace(s, i+1)

		value, next, err := parseYAMLFlowValue(s, i)
		if err != nil {
			return nil, 0, err
		}
		result[key] = value
		i = skipYAMLFlowSpace(s, next)

		if i >= len(s) {
			return nil, 0, fmt.Errorf("unterminated flow mapping")
		}
		switch s[i] {
		case ',':
			i = skipYAMLFlowSpace(s, i+1)
			if i < len(s) && s[i] == '}' {
				return result, i + 1, nil
			}
		case '}':
			return result, i + 1, nil
		default:
			return nil, 0, fmt.Errorf("expected ',' or '}' in flow mapping")
		}
	}
}

func parseYAMLFlowValue(s string, start int) (any, int, error) {
	i := skipYAMLFlowSpace(s, start)
	if i >= len(s) {
		return nil, 0, fmt.Errorf("unexpected end of flow value")
	}

	switch s[i] {
	case '[':
		return parseYAMLFlowSequence(s, i)
	case '{':
		return parseYAMLFlowMapping(s, i)
	case '"', '\'':
		value, next, err := parseYAMLQuoted(s, i)
		if err != nil {
			return nil, 0, err
		}
		return value, next, nil
	}

	valueStart := i
	depth := 0
	for i < len(s) {
		c := s[i]
		if depth == 0 && (c == ',' || c == ']' || c == '}') {
			break
		}
		switch c {
		case '[', '{':
			depth++
		case ']', '}':
			if depth > 0 {
				depth--
			}
		}
		i++
	}

	raw := strings.TrimSpace(s[valueStart:i])
	if raw == "" {
		return nil, 0, fmt.Errorf("empty flow value")
	}
	if isYAMLNull(raw) {
		return nil, i, nil
	}
	return raw, i, nil
}

func parseYAMLQuoted(s string, start int) (string, int, error) {
	if start >= len(s) || (s[start] != '"' && s[start] != '\'') {
		return "", 0, fmt.Errorf("expected quoted string")
	}

	quote := s[start]
	var builder strings.Builder

	for i := start + 1; i < len(s); i++ {
		c := s[i]

		if quote == '\'' {
			if c == '\'' {
				if i+1 < len(s) && s[i+1] == '\'' {
					builder.WriteByte('\'')
					i++
					continue
				}
				return builder.String(), i + 1, nil
			}
			builder.WriteByte(c)
			continue
		}

		if c == '"' {
			return builder.String(), i + 1, nil
		}
		if c != '\\' {
			builder.WriteByte(c)
			continue
		}

		if i+1 >= len(s) {
			return "", 0, fmt.Errorf("unterminated escape sequence")
		}
		i++
		switch esc := s[i]; esc {
		case '0':
			builder.WriteByte(0)
		case 'a':
			builder.WriteByte('\a')
		case 'b':
			builder.WriteByte('\b')
		case 't':
			builder.WriteByte('\t')
		case 'n':
			builder.WriteByte('\n')
		case 'v':
			builder.WriteByte('\v')
		case 'f':
			builder.WriteByte('\f')
		case 'r':
			builder.WriteByte('\r')
		case 'e':
			builder.WriteByte(0x1b)
		case ' ':
			builder.WriteByte(' ')
		case '"':
			builder.WriteByte('"')
		case '/':
			builder.WriteByte('/')
		case '\\':
			builder.WriteByte('\\')
		case 'N':
			builder.WriteRune('\u0085')
		case '_':
			builder.WriteRune('\u00a0')
		case 'L':
			builder.WriteRune('\u2028')
		case 'P':
			builder.WriteRune('\u2029')
		case 'x':
			r, next, err := parseYAMLHexEscape(s, i+1, 2)
			if err != nil {
				return "", 0, err
			}
			builder.WriteRune(r)
			i = next - 1
		case 'u':
			r, next, err := parseYAMLHexEscape(s, i+1, 4)
			if err != nil {
				return "", 0, err
			}
			builder.WriteRune(r)
			i = next - 1
		case 'U':
			r, next, err := parseYAMLHexEscape(s, i+1, 8)
			if err != nil {
				return "", 0, err
			}
			builder.WriteRune(r)
			i = next - 1
		default:
			return "", 0, fmt.Errorf("invalid escape sequence \\%c", esc)
		}
	}

	return "", 0, fmt.Errorf("unterminated quoted string")
}

func parseYAMLHexEscape(s string, start, length int) (rune, int, error) {
	if start+length > len(s) {
		return 0, 0, fmt.Errorf("incomplete hexadecimal escape sequence")
	}
	value, err := strconv.ParseUint(s[start:start+length], 16, 32)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid hexadecimal escape sequence")
	}
	r := rune(value)
	if !utf8.ValidRune(r) {
		return 0, 0, fmt.Errorf("invalid unicode code point")
	}
	return r, start + length, nil
}

func skipYAMLFlowSpace(s string, start int) int {
	for start < len(s) && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	return start
}

// ==================== 配置解码 ====================

func decodeTemplateConfig(root any) (*TemplateConfig, error) {
	cfg := &TemplateConfig{}
	if root == nil {
		return cfg, nil
	}

	m, ok := root.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("top-level YAML value must be a mapping")
	}

	cfg.Name = yamlString(m["name"])
	cfg.Description = yamlString(m["description"])
	cfg.Version = yamlString(m["version"])
	cfg.Author = yamlString(m["author"])

	var err error
	if cfg.Tags, err = yamlStringList(m["tags"]); err != nil {
		return nil, fmt.Errorf("tags: %w", err)
	}
	if cfg.Ignore, err = yamlStringList(m["ignore"]); err != nil {
		return nil, fmt.Errorf("ignore: %w", err)
	}
	if cfg.Delimiters, err = yamlStringList(m["delimiters"]); err != nil {
		return nil, fmt.Errorf("delimiters: %w", err)
	}
	if cfg.Questions, err = decodeQuestions(m["questions"]); err != nil {
		return nil, err
	}
	if cfg.Rename, err = decodeRenameRules(m["rename"]); err != nil {
		return nil, err
	}
	if cfg.Hooks, err = decodeHooks(m["hooks"]); err != nil {
		return nil, err
	}

	return cfg, nil
}

func decodeQuestions(value any) ([]Question, error) {
	if value == nil {
		return nil, nil
	}

	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("questions must be a list")
	}

	questions := make([]Question, 0, len(items))
	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("questions[%d] must be a mapping", i)
		}

		choices, err := yamlStringList(m["choices"])
		if err != nil {
			return nil, fmt.Errorf("questions[%d].choices: %w", i, err)
		}

		questions = append(questions, Question{
			Name:      yamlString(m["name"]),
			Prompt:    yamlString(m["prompt"]),
			Type:      QuestionType(yamlString(m["type"])),
			Default:   yamlString(m["default"]),
			Required:  yamlBool(m["required"]),
			Choices:   choices,
			Validate:  yamlString(m["validate"]),
			Transform: yamlString(m["transform"]),
		})
	}

	return questions, nil
}

func decodeRenameRules(value any) ([]RenameRule, error) {
	if value == nil {
		return nil, nil
	}

	items, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("rename must be a list")
	}

	rules := make([]RenameRule, 0, len(items))
	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("rename[%d] must be a mapping", i)
		}
		rules = append(rules, RenameRule{
			Pattern: yamlString(m["pattern"]),
			Value:   yamlString(m["value"]),
		})
	}

	return rules, nil
}

func decodeHooks(value any) (Hook, error) {
	if value == nil {
		return Hook{}, nil
	}

	m, ok := value.(map[string]any)
	if !ok {
		return Hook{}, fmt.Errorf("hooks must be a mapping")
	}

	before, err := yamlStringList(m["before"])
	if err != nil {
		return Hook{}, fmt.Errorf("hooks.before: %w", err)
	}
	after, err := yamlStringList(m["after"])
	if err != nil {
		return Hook{}, fmt.Errorf("hooks.after: %w", err)
	}

	return Hook{Before: before, After: after}, nil
}

func yamlString(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprint(v)
	}
}

func yamlBool(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	case string:
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "true", "yes", "on", "y", "1":
			return true
		}
	}
	return false
}

func yamlStringList(value any) ([]string, error) {
	switch v := value.(type) {
	case nil:
		return nil, nil
	case []any:
		result := make([]string, 0, len(v))
		for i, item := range v {
			if item == nil {
				result = append(result, "")
				continue
			}
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("item %d must be a string", i)
			}
			result = append(result, s)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("expected a list of strings")
	}
}
