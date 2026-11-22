package ignore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

type rule struct {
	pattern       string
	matcher       *regexp.Regexp
	negated       bool
	directoryOnly bool
	basenameOnly  bool
	anchored      bool
}

type Matcher struct {
	rules []rule
}

func Matches(pattern, name string) bool {
	matcher, err := Parse(strings.NewReader(pattern))
	if err != nil {
		return false
	}
	return matcher.Ignored(name, false)
}

func Parse(source io.Reader) (*Matcher, error) {
	matcher := &Matcher{}
	scanner := bufio.NewScanner(source)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		pattern := strings.TrimSpace(scanner.Text())
		if pattern == "" || strings.HasPrefix(pattern, "#") {
			continue
		}
		entry := rule{}
		if strings.HasPrefix(pattern, "!") {
			entry.negated = true
			pattern = strings.TrimPrefix(pattern, "!")
		}
		if strings.HasPrefix(pattern, "/") {
			entry.anchored = true
			pattern = strings.TrimPrefix(pattern, "/")
		}
		if strings.HasSuffix(pattern, "/") {
			entry.directoryOnly = true
			pattern = strings.TrimSuffix(pattern, "/")
		}
		if pattern == "" {
			return nil, fmt.Errorf("line %d: ignore pattern cannot be empty", lineNumber)
		}
		entry.pattern = pattern
		entry.basenameOnly = !strings.Contains(pattern, "/")
		expression, err := globExpression(pattern)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid ignore pattern %q: %w", lineNumber, pattern, err)
		}
		entry.matcher = expression
		matcher.rules = append(matcher.rules, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read ignore patterns: %w", err)
	}
	return matcher, nil
}

func Load(path string) (*Matcher, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return defaultMatcher(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("open ignore file %s: %w", path, err)
	}
	defer file.Close()
	custom, err := Parse(file)
	if err != nil {
		return nil, fmt.Errorf("parse ignore file %s: %w", path, err)
	}
	matcher := defaultMatcher()
	matcher.rules = append(matcher.rules, custom.rules...)
	return matcher, nil
}

func defaultMatcher() *Matcher {
	matcher, _ := Parse(strings.NewReader(`.git/
.bitpack/
node_modules/
.env
.env.*
*.secret
`))
	return matcher
}

func (m *Matcher) Ignored(name string, isDirectory bool) bool {
	if m == nil {
		return false
	}
	clean := strings.TrimPrefix(path.Clean(strings.ReplaceAll(name, "\\", "/")), "./")
	if clean == "." || strings.HasPrefix(clean, "../") || path.IsAbs(clean) {
		return false
	}
	parts := strings.Split(clean, "/")
	ignored := false
	for _, entry := range m.rules {
		matched := false
		if entry.basenameOnly {
			start := 0
			end := len(parts)
			if entry.anchored {
				end = 1
			}
			for index := start; index < end; index++ {
				isFinal := index == len(parts)-1
				if entry.directoryOnly && isFinal && !isDirectory {
					continue
				}
				if entry.matcher.MatchString(parts[index]) {
					matched = true
					break
				}
			}
		} else {
			for index := range parts {
				isFinal := index == len(parts)-1
				if entry.directoryOnly && isFinal && !isDirectory {
					continue
				}
				candidate := strings.Join(parts[:index+1], "/")
				if entry.matcher.MatchString(candidate) {
					matched = true
					break
				}
			}
		}
		if matched {
			ignored = !entry.negated
		}
	}
	return ignored
}

func globExpression(pattern string) (*regexp.Regexp, error) {
	var expression strings.Builder
	expression.WriteByte('^')
	for index := 0; index < len(pattern); index++ {
		switch pattern[index] {
		case '*':
			if index+1 < len(pattern) && pattern[index+1] == '*' {
				index++
				if index+1 < len(pattern) && pattern[index+1] == '/' {
					expression.WriteString("(?:.*/)?")
					index++
				} else {
					expression.WriteString(".*")
				}
			} else {
				expression.WriteString("[^/]*")
			}
		case '?':
			expression.WriteString("[^/]")
		default:
			expression.WriteString(regexp.QuoteMeta(string(pattern[index])))
		}
	}
	expression.WriteByte('$')
	return regexp.Compile(expression.String())
}
