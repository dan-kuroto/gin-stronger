package goprint

import "strings"

// `isSeparate`: When there is no line break, whether a space needs to be added.
func appendIndent(sb *strings.Builder, currIndent int, indents []int, isSeparate bool) {
	if currIndent > 0 {
		sb.WriteString("\n")
		for _, indent := range indents {
			if indent > 0 {
				sb.WriteString("|")
				sb.WriteString(strings.Repeat(" ", indent-1))
			}
		}
	} else {
		if isSeparate {
			sb.WriteString(" ")
		}
	}
}
