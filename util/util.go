package util

import (
	"bytes"
)

func Bytes2String(b []byte) string {
	pos = bytes.IndexByte(b, 0x0)
	if pos == -1 {
		pos = len(b)
	}
	return string(b[:pos])
}

func Bytes2Strings(b []byte) (ss []string) {
	for _, bs := range bytes.Split(bytes.Trim(b, "\x00"), []byte("\x00")) {
		ss = append(ss, string(bs))
	}
	return
}
