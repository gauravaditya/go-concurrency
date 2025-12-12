package main

import (
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nums := generator(ctx, []int{1, 2, 3, 4, 5})
	squares := squarer(ctx, nums)

	printer(ctx, squares)

	fmt.Println("pipeline complete")
}
