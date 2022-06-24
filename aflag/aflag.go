package aflag

import (
	"strconv"
	"strings"
)

type (
	ArrayFlagInt    []int
	ArrayFlagString []string
)

func (i *ArrayFlagInt) String() string {
	as := make([]string, 0)
	for _, j := range *i {
		as = append(as, strconv.Itoa(j))
	}
	return strings.Join(as, ",")
}

func (i *ArrayFlagInt) Set(value string) error {
	var (
		e error
		j int
	)
	for _, ss := range strings.Split(value, ",") {
		j, e = strconv.Atoi(ss)
		if e == nil {
			*i = append(*i, j)
		} else {
			return e
		}
	}
	return e
}

func (s *ArrayFlagString) String() string {
	return strings.Join(*s, ",")
}

func (s *ArrayFlagString) Set(value string) error {
	*s = append(*s, strings.Split(value, ",")...)
	return nil
}
