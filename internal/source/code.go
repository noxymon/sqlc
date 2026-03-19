package source

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type Edit struct {
	Location int
	Old      string
	New      string
	OldFunc  func(string) int
}

type CommentSyntax struct {
	Dash      bool
	Hash      bool
	SlashStar bool
}

func LineNumber(source string, head int) (int, int) {
	// Calculate the true line and column number for a query, ignoring spaces
	var comment bool
	var loc, line, col int
	for i, char := range source {
		loc += 1
		col += 1
		// TODO: Check bounds
		if char == '-' && source[i+1] == '-' {
			comment = true
		}
		if char == '\n' {
			comment = false
			line += 1
			col = 0
		}
		if loc <= head {
			continue
		}
		if unicode.IsSpace(char) {
			continue
		}
		if comment {
			continue
		}
		break
	}
	return line + 1, col
}

func Pluck(source string, location, length int) (string, error) {
	head := location
	tail := location + length
	return source[head:tail], nil
}

func Mutate(raw string, a []Edit) (string, error) {
	if len(a) == 0 {
		return raw, nil
	}

	sort.Slice(a, func(i, j int) bool { return a[i].Location > a[j].Location })

	s := raw
	for idx, edit := range a {
		start := edit.Location
		if start > len(s) || start < 0 {
			return "", fmt.Errorf("edit start location is out of bounds")
		}
		var oldLen int
		if edit.OldFunc != nil {
			oldLen = edit.OldFunc(s[start:])
		} else {
			oldLen = len(edit.Old)
		}

		stop := edit.Location + oldLen
		if stop > len(s) {
			return "", fmt.Errorf("edit stop location is out of bounds")
		}

		// If this is not the first edit, (applied backwards), check if
		// this edit overlaps the previous one (and is therefore a developer error)
		if idx != 0 {
			prevEdit := a[idx-1]
			if prevEdit.Location < edit.Location+oldLen {
				return "", fmt.Errorf("2 edits overlap")
			}
		}

		s = s[:start] + edit.New + s[stop:]
	}
	return s, nil
}

func StripComments(sql string) (string, []string, error) {
	s := bufio.NewScanner(strings.NewReader(strings.TrimSpace(sql)))
	var lines, comments []string
	var inBlock bool
	for s.Scan() {
		t := s.Text()
		trimmed := strings.TrimSpace(t)

		if inBlock {
			if i := strings.Index(t, "*/"); i != -1 {
				comments = append(comments, t[:i])
				inBlock = false
			} else {
				comments = append(comments, t)
			}
			continue
		}

		if strings.HasPrefix(trimmed, "-- name:") {
			continue
		}
		if strings.HasPrefix(trimmed, "/* name:") && strings.HasSuffix(trimmed, "*/") {
			continue
		}
		if strings.HasPrefix(trimmed, "# name:") {
			continue
		}

		if strings.HasPrefix(trimmed, "--") {
			comments = append(comments, strings.TrimPrefix(trimmed, "--"))
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			comments = append(comments, strings.TrimPrefix(trimmed, "#"))
			continue
		}
		if strings.HasPrefix(trimmed, "/*") {
			if i := strings.Index(t, "*/"); i != -1 {
				// Single line block
				content := t[strings.Index(t, "/*")+2 : i]
				// Check if it's a name annotation (even if not at start of line)
				if !strings.HasPrefix(strings.TrimSpace(content), "name:") {
					comments = append(comments, content)
				}
			} else {
				inBlock = true
				content := t[strings.Index(t, "/*")+2:]
				if !strings.HasPrefix(strings.TrimSpace(content), "name:") {
					comments = append(comments, content)
				}
			}
			continue
		}
		lines = append(lines, t)
	}
	return strings.Join(lines, "\n"), comments, s.Err()
}

func CleanedComments(rawSQL string, cs CommentSyntax) ([]string, error) {
	s := bufio.NewScanner(strings.NewReader(strings.TrimSpace(rawSQL)))
	var comments []string
	var inBlock bool
	for s.Scan() {
		line := s.Text()
		trimmed := strings.TrimSpace(line)

		if inBlock {
			if i := strings.Index(line, "*/"); i != -1 {
				comments = append(comments, line[:i])
				inBlock = false
			} else {
				comments = append(comments, line)
			}
			continue
		}

		if cs.Dash && strings.HasPrefix(trimmed, "--") {
			comments = append(comments, strings.TrimPrefix(trimmed, "--"))
			continue
		}
		if cs.Hash && strings.HasPrefix(trimmed, "#") {
			comments = append(comments, strings.TrimPrefix(trimmed, "#"))
			continue
		}
		if cs.SlashStar && strings.HasPrefix(trimmed, "/*") {
			if i := strings.Index(line, "*/"); i != -1 {
				// Single line block
				content := line[strings.Index(line, "/*")+2 : i]
				comments = append(comments, content)
			} else {
				inBlock = true
				comments = append(comments, line[strings.Index(line, "/*")+2:])
			}
			continue
		}
	}
	return comments, s.Err()
}
