package ui

import "golang.org/x/exp/constraints"

type Numeric interface {
	constraints.Float | constraints.Integer
}

func Max[T Numeric](s []T) (m T) {
	for _, v := range s {
		if v > m {
			m = v
		}
	}
	return
}
