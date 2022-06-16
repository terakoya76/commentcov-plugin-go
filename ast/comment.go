package ast

import (
	"strings"
)

// Normalize the given comment text.
func Normalize(str string) string {
	return strings.TrimLeft(str, " ")
}

func IsOnlyNoLintAnnotation(str string) bool {
	rows := strings.Split(str, "\n")

	return len(rows) == 1 && strings.HasPrefix(Normalize(rows[0]), "nolint:") ||
		len(rows) == 2 && strings.HasPrefix(Normalize(rows[0]), "nolint:") && rows[1] == ""
}
