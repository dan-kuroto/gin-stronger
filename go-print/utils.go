package goprint

import "strings"

// When `indent>0`, Add line break and spaces(for `indent*layer`).
// When `indent<=0`, just add a space if `separate`.
func appendIndent(sb *strings.Builder, indent int, layer int, separate bool) {
	if indent > 0 {
		sb.WriteString("\n")
		sb.WriteString(strings.Repeat(" ", indent*layer))
	} else {
		if separate {
			sb.WriteString(" ")
		}
	}
}
