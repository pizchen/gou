package fstdin

import (
	"bufio"
	"os"
)

func FileStdin() (*bufio.Reader, func()) {
	var r *bufio.Reader
	fc := func() {}
	if len(os.Args) > 1 {
		f, e := os.Open(os.Args[1])
		if e == nil {
			r = bufio.NewReader(f)
			fc = func() { _ = f.Close() }
		} else {
			panic(e)
		}
	}
	if r == nil {
		r = bufio.NewReader(os.Stdin)
	}
	return r, fc
}
