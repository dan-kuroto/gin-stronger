package gp

import "strings"

// `isSeparate`: When there is no line break, whether a space needs to be added.
func appendIndent(sb *strings.Builder, currIndent int, indents []int, isSeparate bool, useColor bool) {
	if currIndent > 0 {
		sb.WriteString("\n")
		for i, indent := range indents {
			if indent > 0 {
				appendColoredString(sb, "|", i, useColor, false)
				sb.WriteString(strings.Repeat(" ", indent-1))
			}
		}
	} else {
		if isSeparate {
			sb.WriteString(" ")
		}
	}
}

func appendColoredString(sb *strings.Builder, str string, colorIdx int, useColor bool, bold bool) {
	if useColor {
		if bold {
			boldColors[colorIdx%len(boldColors)].Fprint(sb, str)
		} else {
			normalColors[colorIdx%len(normalColors)].Fprint(sb, str)
		}
	} else {
		sb.WriteString(str)
	}
}
