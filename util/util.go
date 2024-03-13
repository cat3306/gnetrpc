package util

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/panjf2000/gnet/v2"
)

func IsExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func JoinServiceMethod(s, m string) string {
	return s + "@" + m
}
func GetConnId(conn gnet.Conn) string {
	return fmt.Sprintf("%s@%d", conn.RemoteAddr().String(), conn.Fd())
}
func GoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
