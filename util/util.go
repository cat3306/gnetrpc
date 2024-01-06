package util

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/lithammer/shortuuid/v4"
)

func IsExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func JoinServiceMethod(s, m string) string {
	return s + "@" + m
}

func GenConnId(fd int) string {
	return fmt.Sprintf("%s@%d", shortuuid.New(), fd)
}
