package exclude

import (
	"fmt"
	"strings"
)

func NewCommaSeparatedList(s string) (List, error) {
	l := List(strings.Split(s, ","))
	// aproximately linear algorithm started to be non-optimal if length > ~25-30
	if len(l) > 30 {
		return nil, fmt.Errorf("please use non-linear implementation of an exclude list: length of list %d > 30", len(l))
	}
	return l, nil
}

type List []string

func (e List) Check(v string) (ok bool) {
	for _, c := range e {
		if c == v {
			return false
		}
	}
	return true
}
