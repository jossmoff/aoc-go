package utils

import (
	"iter"

	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Integer | constraints.Float](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Sum[T constraints.Integer | constraints.Float](s iter.Seq[T]) T {
	var acc = T(0)
	for v := range s {
		acc += v
	}
	return acc
}

func Any[T any](f func(T) bool, s iter.Seq[T]) bool {
	for v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

func All[T any](f func(T) bool, s iter.Seq[T]) bool {
	return !Any(f, s)
}

func Reduce[Sum, V any](f func(Sum, V) Sum, sum Sum, s iter.Seq[V]) Sum {
	for v := range s {
		sum = f(sum, v)
	}
	return sum
}

func Map[In, Out any](f func(In) Out, s iter.Seq[In]) iter.Seq[Out] {
	return func(yield func(Out) bool) {
		for in := range s {
			if !yield(f(in)) {
				return
			}
		}
	}
}

func Filter[T any](f func(T) bool, s iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}

type Zipped[V1, V2 any] struct {
	V1  V1
	Ok1 bool // whether V1 is present (if not, it will be zero)
	V2  V2
	Ok2 bool // whether V2 is present (if not, it will be zero)
}

func Zip[V1, V2 any](x iter.Seq[V1], y iter.Seq[V2]) iter.Seq[Zipped[V1, V2]] {
	return func(yield func(z Zipped[V1, V2]) bool) {
		next, stop := iter.Pull(y)
		defer stop()
		v2, ok2 := next()
		for v1 := range x {
			if !yield(Zipped[V1, V2]{v1, true, v2, ok2}) {
				return
			}
			v2, ok2 = next()
		}
		var zv1 V1
		for ok2 {
			if !yield(Zipped[V1, V2]{zv1, false, v2, ok2}) {
				return
			}
			v2, ok2 = next()
		}
	}
}

func Tally[T comparable](in []T) map[T]int {
	ret := make(map[T]int)
	for _, v := range in {
		ret[v] += 1
	}

	return ret
}

func Flatten2D[T any](grid [][]T) []T {
	var result []T
	for _, row := range grid {
		result = append(result, row...)
	}
	return result
}
