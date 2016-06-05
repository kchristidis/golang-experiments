package main

import (
	"fmt"
	"sync"
)

// first stage
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// second stage
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// main
func runPipeline() {
	/* c := gen(2, 3)
	out := sq(c)
	fmt.Println(<-out)
	fmt.Println(<-out) */

	/* for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n)
	} */

	in := gen(2, 3)

	// distribute the sq work across two goroutines that both read from in
	c1 := sq(in)
	c2 := sq(in)

	// consume the merged output from c1 and c2
	for n := range merge(c1, c2) {
		fmt.Println(n)
	}
}
