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
