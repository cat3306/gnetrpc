package util

import (
	"unicode"
	"unicode/utf8"
)

func IsExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func JoinServiceMethod(s, m string) string {
	return s + "@" + m
}
