package shared

import (
	"golang.org/x/exp/constraints"
)

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func ContainsWhere[T any](slice []T, predicate func(T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			return true
		}
	}

	return false
}

func ContainsWhereP[T any](slice []T, predicate func(*T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if predicate(&slice[i]) {
			return true
		}
	}

	return false
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Partition[T any](elems []T, predicate func(T, int, []T) bool) (pass []T, fail []T) {
	pass = make([]T, 0)
	fail = make([]T, 0)

	for i, e := range elems {
		if predicate(e, i, elems) {
			pass = append(pass, e)
		} else {
			fail = append(fail, e)
		}
	}

	return pass, fail
}

func FirstOrDefaultP[T any](slice []T, predicate func(*T) bool) (element *T) {
	for i := 0; i < len(slice); i++ {
		if predicate(&slice[i]) {
			return &slice[i]
		}
	}

	return nil
}

func FirstOrDefault[T any](slice []T, predicate func(T) bool) (element T) {
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			return slice[i]
		}
	}

	var zero T
	return zero
}

func Where[T any](slice []T, predicate func(T) bool) []T {
	ret := make([]T, 0)

	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			ret = append(ret, slice[i])
		}
	}

	return ret
}

func SliceToMap[T any, K comparable](slice []T, keyFunc func(*T) K) map[K]T {
	m := make(map[K]T)
	for _, v := range slice {
		m[keyFunc(&v)] = v
	}
	return m
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
