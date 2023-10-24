package util

import "unsafe"

func BytesToString(raw []byte) string {
	return *(*string)(unsafe.Pointer(&raw))
}
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		Cap int
	}{s, len(s)}))
}
