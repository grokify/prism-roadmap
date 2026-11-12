// Package mdshift moves Markdown headings down a number of levels, so a
// rendered document can be embedded as a section of a larger one.
package mdshift

import "strings"

// MaxLevel is the deepest Markdown heading level.
const MaxLevel = 6

// Shift pushes every ATX heading ("#" to "######" followed by a space, or
// alone on its line) down n levels, capped at level 6. Lines inside fenced
// code blocks (``` or ~~~) are left untouched, as is everything else. A
// non-positive n returns md unchanged.
func Shift(md string, n int) string {
	if n <= 0 || md == "" {
		return md
	}
	lines := strings.Split(md, "\n")
	var fence string // the opening fence while inside a code block
	for i, line := range lines {
		if fence != "" {
			if closesFence(line, fence) {
				fence = ""
			}
			continue
		}
		if f := openingFence(line); f != "" {
			fence = f
			continue
		}
		lines[i] = shiftHeading(line, n)
	}
	return strings.Join(lines, "\n")
}

// shiftHeading returns line with its heading level raised by n, or line
// unchanged if it is not an ATX heading.
func shiftHeading(line string, n int) string {
	indent := leadingSpaces(line)
	if indent > 3 {
		return line
	}
	rest := line[indent:]
	level := 0
	for level < len(rest) && rest[level] == '#' {
		level++
	}
	if level == 0 || level > MaxLevel {
		return line
	}
	if level < len(rest) && rest[level] != ' ' && rest[level] != '\t' {
		return line
	}
	return line[:indent] + strings.Repeat("#", min(level+n, MaxLevel)) + rest[level:]
}

// openingFence returns the fence (e.g. "```" or "~~~~") that line opens, or
// "" if line does not open a fenced code block.
func openingFence(line string) string {
	indent := leadingSpaces(line)
	if indent > 3 {
		return ""
	}
	rest := line[indent:]
	if rest == "" || (rest[0] != '`' && rest[0] != '~') {
		return ""
	}
	c := rest[0]
	run := 0
	for run < len(rest) && rest[run] == c {
		run++
	}
	if run < 3 {
		return ""
	}
	if c == '`' && strings.Contains(rest[run:], "`") {
		return "" // a backtick info string cannot contain backticks
	}
	return rest[:run]
}

// closesFence reports whether line closes a code block opened by fence: the
// same character, at least as many of it, and nothing but spaces after.
func closesFence(line, fence string) bool {
	indent := leadingSpaces(line)
	if indent > 3 {
		return false
	}
	rest := line[indent:]
	c := fence[0]
	run := 0
	for run < len(rest) && rest[run] == c {
		run++
	}
	return run >= len(fence) && strings.TrimSpace(rest[run:]) == ""
}

func leadingSpaces(s string) int {
	n := 0
	for n < len(s) && s[n] == ' ' {
		n++
	}
	return n
}
