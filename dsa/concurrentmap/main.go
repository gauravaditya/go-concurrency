package main

func ParallelMap[T any, R any](
	in []T,
	workers int,
	fn func(T) R,
) []R {
	
}

type input[T, R any] struct {
	v   T
	r   R
	idx int
}
