package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func runShell(input io.Reader, output io.Writer, run func([]string) error) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	showPrompt := false
	if file, ok := input.(*os.File); ok {
		if info, err := file.Stat(); err == nil {
			showPrompt = info.Mode()&os.ModeCharDevice != 0
		}
	}

	fmt.Fprintln(output, "Bitpack Go shell. Type 'help' for commands, 'exit' to quit.")
	for {
		if showPrompt {
			fmt.Fprint(output, "bp> ")
		}
		if !scanner.Scan() {
			return scanner.Err()
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		args, err := splitCommand(line)
		if err != nil {
			fmt.Fprintf(output, "Error: %v\n", err)
			continue
		}
		if len(args) == 0 {
			continue
		}
		if args[0] == "exit" || args[0] == "quit" {
			return nil
		}
		if err := run(args); err != nil {
			fmt.Fprintf(output, "Error: %v\n", err)
		}
	}
}

func splitCommand(line string) ([]string, error) {
	var args []string
	var current strings.Builder
	var quote rune
	escaped := false
	inToken := false

	characters := []rune(line)
	for i, char := range characters {
		if escaped {
			current.WriteRune(char)
			escaped = false
			inToken = true
			continue
		}
		if char == '\\' && quote != '\'' {
			if i+1 < len(characters) {
				next := characters[i+1]
				if next == '\\' || next == '"' || next == '\'' || next == ' ' || next == '\t' {
					escaped = true
					inToken = true
					continue
				}
			}
			current.WriteRune(char)
			inToken = true
			continue
		}
		if quote != 0 {
			if char == quote {
				quote = 0
			} else {
				current.WriteRune(char)
			}
			inToken = true
			continue
		}
		switch char {
		case '\'', '"':
			quote = char
			inToken = true
		case ' ', '\t':
			if inToken {
				args = append(args, current.String())
				current.Reset()
				inToken = false
			}
		default:
			current.WriteRune(char)
			inToken = true
		}
	}
	if escaped {
		current.WriteRune('\\')
	}
	if quote != 0 {
		return nil, fmt.Errorf("unclosed quote")
	}
	if inToken {
		args = append(args, current.String())
	}
	return args, nil
}
