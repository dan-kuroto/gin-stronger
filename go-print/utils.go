package goprint

import "strings"

// When `indent>0`, Add line break and spaces(for `indent*layer`).
// When `indent<=0`, just add a space if `isSeparate`.
func appendIndent(sb *strings.Builder, indent int, layer int, isSeparate bool) {
	if indent > 0 {
		sb.WriteString("\n")
		for i := 0; i < layer; i++ {
			sb.WriteString("|")
			sb.WriteString(strings.Repeat(" ", indent-1))
		}
	} else {
		if isSeparate {
			sb.WriteString(" ")
		}
	}
}

// BUG: 目前如果几个indent不同（除零之外），在不同类型嵌套的时候会错位……
